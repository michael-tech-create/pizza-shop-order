// =============================================================
// auth-guard.js — Session protection for admin-only pages
//
// Include this on EVERY page that should require login
// (admin.html, activity.html, etc.) BEFORE any other script that
// makes API calls, so the redirect happens before the page tries
// to fetch protected data.
//
// Usage in HTML:
//   <script src="./js/auth-guard.js"></script>
//   <script src="./js/logger.js"></script>
//   ... rest of page scripts
// =============================================================

const AuthGuard = (() => {
    const TOKEN_KEY    = "mcpizza_admin_token";
    const USERNAME_KEY = "mcpizza_admin_username";
    const EXPIRY_KEY   = "mcpizza_admin_expiry";
    const LOGIN_PAGE   = "login.html";

    function getToken()    { return localStorage.getItem(TOKEN_KEY); }
    function getUsername() { return localStorage.getItem(USERNAME_KEY); }
    function getExpiry()   { return Number(localStorage.getItem(EXPIRY_KEY) || 0); }

    function isExpired() {
        const exp = getExpiry();
        return !exp || (Date.now() / 1000) >= exp;
    }

    function clearSession() {
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USERNAME_KEY);
        localStorage.removeItem(EXPIRY_KEY);
    }

    function redirectToLogin(reason) {
        clearSession();
        const params = reason ? `?reason=${encodeURIComponent(reason)}` : "";
        window.location.href = `${LOGIN_PAGE}${params}`;
    }

    function logout() {
        clearSession();
        window.location.href = LOGIN_PAGE;
    }

    // ── Synchronous guard — runs immediately on script load ──
    // If there's no token or it's expired, bounce to login before the
    // rest of the page even renders, avoiding a flash of protected content.
    (function enforceAuth() {
        const token = getToken();
        if (!token || isExpired()) {
            redirectToLogin(!token ? "no_session" : "expired");
        }
    })();

    // ── authFetch — drop-in replacement for fetch() that automatically
    //    attaches the Authorization header and handles 401 globally ──
    async function authFetch(url, options = {}) {
        const token = getToken();

        const headers = {
            ...(options.headers || {}),
            "Authorization": `Bearer ${token}`,
        };

        const res = await fetch(url, { ...options, headers });

        if (res.status === 401) {
            // Token rejected by the server (expired/invalid/secret rotated)
            redirectToLogin("session_invalid");
            throw new Error("Session expired");
        }

        return res;
    }

    // ── Warn the user a few minutes before their session expires so they
    //    don't lose unsaved work mid-task (e.g. mid-edit on a pizza form) ──
    function startExpiryWarning(onWarn) {
        const checkInterval = 30_000; // check every 30s
        const warnWindow    = 5 * 60; // warn when < 5 min remain

        setInterval(() => {
            const exp = getExpiry();
            if (!exp) return;
            const secondsLeft = exp - Date.now() / 1000;
            if (secondsLeft > 0 && secondsLeft <= warnWindow) {
                if (typeof onWarn === "function") onWarn(Math.round(secondsLeft));
            }
        }, checkInterval);
    }

    return {
        getToken,
        getUsername,
        isExpired,
        logout,
        authFetch,
        startExpiryWarning,
    };
})();

window.AuthGuard = AuthGuard;