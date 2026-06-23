

const API_BASE = "http://localhost:8080";


function showToast(message, type = "success") {
    const existing = document.getElementById("orderToast");
    if (existing) existing.remove();

    const colours = {
        success: "bg-green-600",
        error:   "bg-red-600",
        info:    "bg-blue-600",
    };

    const toast = document.createElement("div");
    toast.id = "orderToast";
    toast.className = `fixed bottom-6 left-1/2 -translate-x-1/2 z-[999] px-6 py-3 rounded-xl text-white font-semibold shadow-xl ${colours[type] || colours.success} transition-all`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => toast.remove(), 3500);
}


async function checkout() {
    const grouped = getCartSnapshot();

    if (!grouped || grouped.length === 0) {
        showToast("Your cart is empty!", "error");
        return;
    }

    const checkoutBtn = document.getElementById("checkoutBtn");
    if (checkoutBtn) {
        checkoutBtn.disabled = true;
        checkoutBtn.textContent = "Placing Order…";
    }

    try {
        const results = await Promise.all(
            grouped.map(item =>
                fetch(`${API_BASE}/api/orders`, {
                    method:  "POST",
                    headers: { "Content-Type": "application/json" },
                    body:    JSON.stringify({
                        pizza_id: item.id,
                        quantity: item.quantity,
                    }),
                }).then(async res => {
                    const data = await res.json();
                    if (!res.ok) throw new Error(data.error || "Order failed");
                    return data;
                })
            )
        );

        const totalCost = results.reduce((sum, o) => sum + o.total_cost, 0);

        showToast(`Order placed! Total: ${formatNaira(totalCost)}`, "success");
        clearCart();

    } catch (err) {
        console.error("Checkout error:", err);
        showToast(err.message || "Checkout failed. Please try again.", "error");
    } finally {
        if (checkoutBtn) {
            checkoutBtn.disabled = false;
            checkoutBtn.textContent = "Proceed to Checkout";
        }
    }
}

async function loadOrders() {
    try {
        const res  = await fetch(`${API_BASE}/api/orders`);
        if (!res.ok) throw new Error(`Server error: ${res.status}`);
        const orders = await res.json();
        return orders || [];
    } catch (err) {
        console.error("Failed to load orders:", err);
        return [];
    }
}


function renderOrdersTable(orders, targetEl) {
    if (!targetEl) return;

    if (!orders || orders.length === 0) {
        targetEl.innerHTML = `
            <tr>
                <td colspan="5" class="px-6 py-16 text-center text-slate-400 font-medium">
                    No orders yet.
                </td>
            </tr>
        `;
        return;
    }

    const statusStyles = {
        pending:   "bg-yellow-100 text-yellow-800",
        delivered: "bg-green-100 text-green-700",
        cancelled: "bg-red-100 text-red-700",
    };

    targetEl.innerHTML = orders.map(order => `
        <tr class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4 font-semibold text-slate-700">#${order.order_id}</td>
            <td class="px-6 py-4 text-slate-600">${order.pizza_name}</td>
            <td class="px-6 py-4 text-center text-slate-600">${order.quantity}</td>
            <td class="px-6 py-4 text-right font-bold text-slate-800">${formatNaira(order.total_cost)}</td>
            <td class="px-6 py-4 text-center">
                <select
                    onchange="updateOrderStatus(${order.order_id}, this.value)"
                    class="text-xs font-semibold px-3 py-1.5 rounded-lg border border-slate-200 cursor-pointer focus:outline-none focus:ring-2 focus:ring-orange-500 ${statusStyles[order.status] || statusStyles.pending}"
                >
                    <option value="pending"   ${order.status === "pending"   ? "selected" : ""}>Pending</option>
                    <option value="delivered" ${order.status === "delivered" ? "selected" : ""}>Delivered</option>
                    <option value="cancelled" ${order.status === "cancelled" ? "selected" : ""}>Cancelled</option>
                </select>
            </td>
        </tr>
    `).join("");
}

async function updateOrderStatus(orderId, status) {
    try {
        const res = await fetch(`${API_BASE}/api/orders/${orderId}/status`, {
            method:  "PATCH",
            headers: { "Content-Type": "application/json" },
            body:    JSON.stringify({ status }),
        });

        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error || "Update failed");
        }

        showToast(`Order #${orderId} marked as ${status}`, "success");

    } catch (err) {
        console.error("Status update error:", err);
        showToast(err.message || "Failed to update status", "error");
    }
}


window.checkout = checkout;
window.loadOrders = loadOrders;
window.renderOrdersTable = renderOrdersTable;
window.updateOrderStatus = updateOrderStatus;
window.showToast = showToast;