const API_BASE = "http://localhost:8080";

function showToast(message, type = "success") {
    document.getElementById("orderToast")?.remove();

    const palette = {
        success: "bg-green-600",
        error:   "bg-red-600",
        info:    "bg-orange-500",
    };

    const el = document.createElement("div");
    el.id = "orderToast";
    el.className = [
        "fixed bottom-6 left-1/2 -translate-x-1/2 z-[9999]",
        "px-6 py-3 rounded-xl text-white text-sm font-semibold shadow-2xl",
        "transition-all duration-300",
        palette[type] ?? palette.success,
    ].join(" ");
    el.textContent = message;
    document.body.appendChild(el);
    setTimeout(() => el.remove(), 3200);
}

function openCheckout() {
    const snap = getCartSnapshot();
    if (!snap || snap.length === 0) {
        showToast("Your cart is empty!", "error");
        return;
    }
    resetCheckoutFlow();
    document.getElementById("checkoutModal")?.classList.remove("hidden");
}

function closeCheckout() {
    document.getElementById("checkoutModal")?.classList.add("hidden");
}

function goToPayment() {
    const name    = document.getElementById("co-name")?.value.trim();
    const phone   = document.getElementById("co-phone")?.value.trim();
    const address = document.getElementById("co-address")?.value.trim();

    if (!name || !phone || !address) {
        showToast("Please fill in all delivery details", "error");
        return;
    }

    if (!/^(\+234|0)[789]\d{9}$/.test(phone)) {
        showToast("Enter a valid Nigerian phone number", "error");
        return;
    }

    const snap = getCartSnapshot();
    const total = snap.reduce((s, i) => s + i.price * i.quantity, 0);

    const summaryEl = document.getElementById("orderSummary");
    if (summaryEl) {
        summaryEl.innerHTML = `
            <div class="space-y-2 mb-4">
                ${snap.map(item => `
                    <div class="flex justify-between text-sm py-2 border-b border-gray-100">
                        <span class="text-gray-700">${item.name} <span class="text-gray-400">× ${item.quantity}</span></span>
                        <span class="font-bold text-gray-900">${formatNaira(item.price * item.quantity)}</span>
                    </div>
                `).join("")}
            </div>
            <div class="flex justify-between font-black text-base pt-2">
                <span>Total</span>
                <span class="text-orange-600">${formatNaira(total)}</span>
            </div>
            <div class="mt-4 text-xs text-gray-500 space-y-1 bg-gray-50 rounded-xl p-3">
                <p>📍 <span class="font-medium">${address}</span></p>
                <p>📞 <span class="font-medium">${phone}</span></p>
                <p>👤 <span class="font-medium">${name}</span></p>
            </div>
        `;
    }

    document.getElementById("stepAddress")?.classList.add("hidden");
    document.getElementById("stepConfirm")?.classList.remove("hidden");
}

function backToAddress() {
    document.getElementById("stepConfirm")?.classList.add("hidden");
    document.getElementById("stepAddress")?.classList.remove("hidden");
}

async function startCheckout() {
    const snap = getCartSnapshot();
    if (!snap || snap.length === 0) {
        showToast("Cart is empty!", "error");
        return;
    }

    const name = document.getElementById("co-name")?.value.trim();
    const address = document.getElementById("co-address")?.value.trim();
    const phone = document.getElementById("co-phone")?.value.trim();

    // Map and sanitize layout variables strictly into integers
    const formattedItems = snap.map(item => {
        const pId = Number(item.id) || Number(item.pizza_id) || 0;
        const qty = Number(item.quantity) || 1;
        return { pizza_id: pId, quantity: qty };
    });

    // Client-Side Safeguard against 400 Bad Requests
    if (formattedItems.some(item => item.pizza_id <= 0)) {
        showToast("Cart contains item with an invalid ID. Please clear cart and try again.", "error");
        return;
    }

    const btn = document.getElementById("confirmOrderBtn");
    if (btn) { btn.disabled = true; btn.textContent = "Placing order…"; }

    try {
       const payload = {
        customer_name: document.getElementById("co-name").value,
        phone: document.getElementById("co-phone").value,
        address: document.getElementById("co-address").value,
        items: snap.map(i => ({ pizza_id: i.id, quantity: i.quantity }))
    };

    const res = await fetch(`${API_BASE}/api/orders`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    });

        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error || "Order failed");
        }

        showToast(`Order placed successfully!`, "success");
        clearCart();
        resetCheckoutFlow();
        
        setTimeout(() => {
            window.location.href = "index.html";
        }, 2000);

    } catch (err) {
        console.error("Checkout error:", err);
        showToast(err.message || "Checkout failed. Try again.", "error");
    } finally {
        if (btn) { btn.disabled = false; btn.textContent = "Place Order"; }
    }
}

function resetCheckoutFlow() {
    ["co-name", "co-phone", "co-address"].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.value = "";
    });
    document.getElementById("stepAddress")?.classList.remove("hidden");
    document.getElementById("stepConfirm")?.classList.add("hidden");
    const s = document.getElementById("orderSummary");
    if (s) s.innerHTML = "";
}

async function loadOrders() {
    try {
        const res = await fetch(`${API_BASE}/api/orders`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return await res.json() || [];
    } catch (err) {
        console.error("loadOrders:", err);
        return [];
    }
}

function renderOrdersTable(orders, tbody) {
    if (!tbody) return;

    if (!orders || orders.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="5" class="px-6 py-16 text-center text-slate-400">
                    <svg class="w-12 h-12 mx-auto mb-3 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                            d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
                    </svg>
                    <span class="font-medium">No orders yet</span>
                </td>
            </tr>`;
        return;
    }

    const badge = {
        pending:   "bg-yellow-100 text-yellow-800 border-yellow-200",
        preparing: "bg-blue-100 text-blue-800 border-blue-200",
        delivered: "bg-green-100 text-green-700 border-green-200",
        cancelled: "bg-red-100 text-red-700 border-red-200",
    };
    const emoji = { pending: "🟠", preparing: "🔵", delivered: "🟢", cancelled: "🔴" };

    tbody.innerHTML = orders.map(o => {
        const status = (o.status || "pending").toLowerCase();
        const cost   = Number(o.total_cost || 0);
        return `
            <tr class="hover:bg-amber-50/40 transition-colors">
                <td class="px-6 py-4 font-semibold text-slate-700">#${o.order_id}</td>
                <td class="px-6 py-4 text-slate-600">${o.pizza_name || "—"}</td>
                <td class="px-6 py-4 text-center font-medium text-slate-600">${o.quantity}</td>
                <td class="px-6 py-4 text-right font-bold text-slate-800">${formatNaira(cost)}</td>
                <td class="px-6 py-4 text-center">
                    <select
                        onchange="updateOrderStatus(${o.order_id}, this.value, this)"
                        class="text-xs font-semibold px-2 py-1.5 rounded-lg border cursor-pointer focus:outline-none focus:ring-2 focus:ring-orange-400 ${badge[status] ?? badge.pending}"
                    >
                        ${["pending","preparing","delivered","cancelled"].map(s => `
                            <option value="${s}" ${status === s ? "selected" : ""}>${emoji[s]} ${s.charAt(0).toUpperCase() + s.slice(1)}</option>
                        `).join("")}
                    </select>
                </td>
            </tr>`;
    }).join("");
}

async function updateOrderStatus(orderId, status, selectEl) {
    try {
        const res = await fetch(`${API_BASE}/api/orders/${orderId}/status`, {
            method:  "PATCH",
            headers: { "Content-Type": "application/json" },
            body:    JSON.stringify({ status }),
        });
        if (!res.ok) {
            const d = await res.json();
            throw new Error(d.error || "Update failed");
        }
        showToast(`Order #${orderId} → ${status}`, "success");

        if (typeof refreshStats === "function") refreshStats();
    } catch (err) {
        showToast(err.message, "error");
        if (selectEl) selectEl.value = selectEl.dataset.prev ?? "pending";
    }

    if (selectEl) selectEl.dataset.prev = status;
}

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("checkoutBtn")?.addEventListener("click", openCheckout);
    document.getElementById("nextStepBtn")?.addEventListener("click", goToPayment);
    document.getElementById("backToAddressBtn")?.addEventListener("click", backToAddress);
    document.getElementById("confirmOrderBtn")?.addEventListener("click", startCheckout);
});

window.openCheckout      = openCheckout;
window.closeCheckout     = closeCheckout;
window.goToPayment       = goToPayment;
window.backToAddress     = backToAddress;
window.startCheckout     = startCheckout;
window.updateOrderStatus = updateOrderStatus;
window.showToast         = showToast;
window.loadOrders        = loadOrders;
window.renderOrdersTable = renderOrdersTable;