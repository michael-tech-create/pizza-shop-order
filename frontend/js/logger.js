// =============================================================
// logger.js  —  Client-side activity log
// Intercepts all fetch() calls so every API interaction is
// recorded automatically. Pages can also call AppLogger.log()
// directly for UI events.
//
// Storage: localStorage key "mcpizza_activity_log"
// Max entries kept: 500 (oldest pruned automatically)
// =============================================================

const AppLogger = (() => {
    const STORAGE_KEY = "mcpizza_activity_log";
    const MAX_ENTRIES  = 500;

    // ── Level definitions ──────────────────────────────────
    const LEVELS = {
        info:    { label: "INFO",    color: "#3b82f6", bg: "#eff6ff", border: "#bfdbfe" },
        success: { label: "SUCCESS", color: "#16a34a", bg: "#f0fdf4", border: "#bbf7d0" },
        warning: { label: "WARNING", color: "#d97706", bg: "#fffbeb", border: "#fde68a" },
        error:   { label: "ERROR",   color: "#dc2626", bg: "#fef2f2", border: "#fecaca" },
    };

    // ── Category tags ──────────────────────────────────────
    const CATEGORIES = {
        ORDER:  "🛒 Order",
        PIZZA:  " Pizza",
        IMAGE:  " Image",
        CART:   "🛍️  Cart",
        AUTH:   "🔐 Auth",
        SYSTEM: "⚙️  System",
        UI:     "🖥️  UI",
    };

    // ── Internal helpers ───────────────────────────────────
    function load() {
        try {
            return JSON.parse(localStorage.getItem(STORAGE_KEY)) || [];
        } catch {
            return [];
        }
    }

    function save(entries) {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(entries));
        } catch {
            // localStorage full — prune and retry
            const pruned = entries.slice(-Math.floor(MAX_ENTRIES / 2));
            localStorage.setItem(STORAGE_KEY, JSON.stringify(pruned));
        }
    }

    // ── Public API ─────────────────────────────────────────
    function log(message, level = "info", category = "SYSTEM", meta = {}) {
        const entries = load();

        const entry = {
            id:        crypto.randomUUID?.() ?? Date.now().toString(36) + Math.random().toString(36).slice(2),
            timestamp: new Date().toISOString(),
            level,
            category:  CATEGORIES[category] ?? category,
            message,
            meta,        // arbitrary extra data (url, status, duration…)
        };

        entries.unshift(entry);             // newest first
        if (entries.length > MAX_ENTRIES) entries.splice(MAX_ENTRIES);
        save(entries);

        // Also emit to browser console with colour
        const lvl = LEVELS[level] ?? LEVELS.info;
        console[level === "error" ? "error" : level === "warning" ? "warn" : "log"](
            `%c[${lvl.label}] [${entry.category}] ${message}`,
            `color:${lvl.color};font-weight:600`,
            meta
        );

        // Notify any live listener (activity.html live feed)
        window.dispatchEvent(new CustomEvent("applog", { detail: entry }));

        return entry;
    }

    function getAll()                  { return load(); }
    function clear()                   { save([]); }
    function getLevels()               { return LEVELS; }
    function getCategories()           { return CATEGORIES; }

    function exportCSV() {
        const entries = load();
        if (!entries.length) return;

        const header = "ID,Timestamp,Level,Category,Message,Meta\n";
        const rows = entries.map(e =>
            [
                e.id,
                e.timestamp,
                e.level,
                `"${e.category.replace(/"/g, '""')}"`,
                `"${e.message.replace(/"/g, '""')}"`,
                `"${JSON.stringify(e.meta).replace(/"/g, '""')}"`,
            ].join(",")
        ).join("\n");

        const blob = new Blob([header + rows], { type: "text/csv" });
        const url  = URL.createObjectURL(blob);
        const a    = document.createElement("a");
        a.href     = url;
        a.download = `mcpizza-activity-${new Date().toISOString().slice(0, 10)}.csv`;
        a.click();
        URL.revokeObjectURL(url);

        log("Activity log exported to CSV", "info", "SYSTEM");
    }


    const _origFetch = window.fetch.bind(window);

    window.fetch = async function (input, init = {}) {
        const url    = typeof input === "string" ? input : input.url;
        const method = (init.method || "GET").toUpperCase();
        const start  = performance.now();

 
        let category = "SYSTEM";
        if (/\/api\/orders/.test(url))              category = "ORDER";
        else if (/\/api\/pizzas\/\d+\/images/.test(url)) category = "IMAGE";
        else if (/\/api\/pizzas/.test(url))         category = "PIZZA";
        else if (/\/menu/.test(url))                category = "PIZZA";
        else if (/\/api\/admin/.test(url))          category = "SYSTEM";

        // Friendly message builders per route
        function buildMessage(status) {
            const ok = status >= 200 && status < 300;
            // Order creation
            if (method === "POST" && /\/api\/orders$/.test(url)) {
                let qty = "?", pid = "?";
                try { const b = JSON.parse(init.body); qty = b.quantity; pid = b.pizza_id; } catch {}
                return ok
                    ? `Order placed — pizza #${pid} × ${qty}`
                    : `Order failed — pizza #${pid} × ${qty} (HTTP ${status})`;
            }
            // Order status update
            if (method === "PATCH" && /\/api\/orders\/\d+\/status/.test(url)) {
                const id  = url.match(/\/api\/orders\/(\d+)/)?.[1] ?? "?";
                let st = "?";
                try { st = JSON.parse(init.body).status; } catch {}
                return ok
                    ? `Order #${id} status → ${st}`
                    : `Failed to update order #${id} status (HTTP ${status})`;
            }
            // Pizza create
            if (method === "POST" && /\/api\/pizzas$/.test(url)) {
                let name = "?";
                try { name = JSON.parse(init.body).name; } catch {}
                return ok ? `Pizza created — "${name}"` : `Failed to create pizza "${name}" (HTTP ${status})`;
            }
            // Pizza update
            if (method === "PUT" && /\/api\/pizzas\/\d+$/.test(url)) {
                const id = url.match(/\/api\/pizzas\/(\d+)/)?.[1] ?? "?";
                return ok ? `Pizza #${id} updated` : `Failed to update pizza #${id} (HTTP ${status})`;
            }
            // Pizza delete
            if (method === "DELETE" && /\/api\/pizzas\/\d+$/.test(url)) {
                const id = url.match(/\/api\/pizzas\/(\d+)/)?.[1] ?? "?";
                return ok ? `Pizza #${id} deleted` : `Failed to delete pizza #${id} (HTTP ${status})`;
            }
            // Image upload
            if (method === "POST" && /\/api\/pizzas\/\d+\/images/.test(url)) {
                const id = url.match(/\/api\/pizzas\/(\d+)/)?.[1] ?? "?";
                return ok ? `Image(s) uploaded for pizza #${id}` : `Image upload failed for pizza #${id} (HTTP ${status})`;
            }
            // Image delete
            if (method === "DELETE" && /\/api\/images\/\d+/.test(url)) {
                const id = url.match(/\/api\/images\/(\d+)/)?.[1] ?? "?";
                return ok ? `Image #${id} deleted` : `Failed to delete image #${id} (HTTP ${status})`;
            }
           
            if (method === "GET" && /\/menu/.test(url)) {
                return ok ? `Menu loaded` : `Menu fetch failed (HTTP ${status})`;
            }
           
            return `${method} ${url.replace("http://localhost:8080", "")} → ${status}`;
        }

        try {
            const res      = await _origFetch(input, init);
            const duration = Math.round(performance.now() - start);
            const ok       = res.status >= 200 && res.status < 300;
            const level    = ok ? (method === "GET" ? "info" : "success") : "error";

            log(buildMessage(res.status), level, category, {
                url, method, status: res.status, duration_ms: duration,
            });

            return res;
        } catch (err) {
            const duration = Math.round(performance.now() - start);
            log(`Network error — ${method} ${url.replace("http://localhost:8080", "")} — ${err.message}`,
                "error", category, { url, method, error: err.message, duration_ms: duration });
            throw err;
        }
    };

    // ── Boot message ───────────────────────────────────────
    log("McMichael Pizza admin session started", "info", "SYSTEM", {
        userAgent: navigator.userAgent.slice(0, 80),
        page: window.location.pathname,
    });

    return { log, getAll, clear, exportCSV, getLevels, getCategories, LEVELS, CATEGORIES };
})();

window.AppLogger = AppLogger;