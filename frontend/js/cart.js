// Load from localStorage or start empty
let cart = JSON.parse(localStorage.getItem("cart")) || [];

// Save to localStorage whenever cart changes
function saveCart() {
    localStorage.setItem("cart", JSON.stringify(cart));
}

function formatNaira(amount) {
    return "₦" + Number(amount).toLocaleString("en-NG");
}

function setupCartButtons() {
    document.querySelectorAll(".add-to-cart").forEach(button => {
        button.addEventListener("click", () => {
            // Support both data-id and data-pizza-id safely
            const itemID = Number(button.dataset.id) || Number(button.dataset.pizzaId) || 0;
            
            addToCart({
                id:    itemID,
                name:  button.dataset.name || "Delicious Pizza",
                price: Number(button.dataset.price) || 0,
                image: button.dataset.image || "",
            });
        });
    });
}

function addToCart(item) {
    if (!item.id || item.id <= 0) {
        console.error("Cannot add item: Invalid or missing ID", item);
        return;
    }
    cart.push(item);
    saveCart();
    updateCartCount();
    renderCart();
}

function groupCartItems(items) {
    const grouped = new Map();
    for (const item of items) {
        // Fallback check for alternate ID keys
        const id = Number(item.id) || Number(item.pizza_id) || 0;
        if (id <= 0) continue; // Skip corrupted entries

        const existing = grouped.get(id);
        if (existing) {
            existing.quantity += (Number(item.quantity) || 1);
        } else {
            grouped.set(id, { 
                ...item, 
                id: id, 
                quantity: Number(item.quantity) || 1 
            });
        }
    }
    return Array.from(grouped.values());
}

function updateCartCount() {
    const cartCountEl = document.getElementById("cartCount");
    if (cartCountEl) cartCountEl.textContent = cart.length;
}

function renderCart() {
    const cartItemsEl = document.getElementById("cartItems");
    const cartTotalEl = document.getElementById("cartTotal");
    
    const grouped = groupCartItems(cart);
    let total = 0;

    const html = grouped.map(item => {
        const itemTotal = item.price * item.quantity;
        total += itemTotal;

        return `
            <div class="py-4 flex justify-between items-center border-b border-gray-100 last:border-b-0 hover:bg-gray-50/50 px-2 rounded-xl transition-colors">
                <div class="flex items-center gap-3">
                    ${item.image ? `<img src="${item.image}" class="w-12 h-12 rounded-lg object-cover" alt="${item.name}">` : ""}
                    <div>
                        <p class="font-bold text-gray-800 text-sm flex items-center gap-1.5">
                            ${item.name}
                            ${item.quantity > 1 ? `<span class="bg-orange-50 text-orange-600 font-semibold text-xs px-2 py-0.5 rounded-md">x${item.quantity}</span>` : ""}
                        </p>
                        <p class="text-sm font-black text-gray-900 mt-0.5">${formatNaira(itemTotal)}</p>
                    </div>
                </div>

                <div class="flex items-center gap-2.5 bg-gray-100/80 p-1 rounded-xl border border-gray-200/40">
                    <button
                        onclick="decreaseQuantity(${item.id})"
                        class="bg-white hover:bg-orange-50 hover:text-orange-600 text-gray-600 font-bold w-7 h-7 rounded-lg shadow-sm flex items-center justify-center active:scale-90 transition-all text-xs"
                        aria-label="Decrease quantity"
                    >—</button>

                    <span class="text-xs font-bold text-gray-800 px-1 min-w-[12px] text-center">${item.quantity}</span>

                    <button
                        onclick="increaseQuantity(${item.id})"
                        class="bg-white hover:bg-orange-50 hover:text-orange-600 text-gray-600 font-bold w-7 h-7 rounded-lg shadow-sm flex items-center justify-center active:scale-90 transition-all text-sm"
                        aria-label="Increase quantity"
                    >+</button>
                </div>
            </div>
        `;
    }).join("");

    if (cartItemsEl) {
        cartItemsEl.innerHTML = html || `
            <div class="text-center py-12 flex flex-col items-center justify-center">
                <span class="text-4xl mb-2 opacity-60">0</span>
                <p class="text-gray-400 font-medium text-sm">Your cart is empty.</p>
            </div>
        `;
    }

    if (cartTotalEl) cartTotalEl.textContent = total.toLocaleString("en-NG");
}

function increaseQuantity(id) {
    const targetId = Number(id);
    const item = cart.find(i => (Number(i.id) || Number(i.pizza_id)) === targetId);
    if (item) {
        cart.push({ id: targetId, name: item.name, price: item.price, image: item.image });
        saveCart();
    }
    updateCartCount();
    renderCart();
}

function decreaseQuantity(id) {
    const targetId = Number(id);
    const index = cart.findIndex(i => (Number(i.id) || Number(i.pizza_id)) === targetId);
    if (index !== -1) {
        cart.splice(index, 1);
        saveCart();
    }
    updateCartCount();
    renderCart();
}

function clearCart() {
    cart = [];
    saveCart();
    updateCartCount();
    renderCart();
}

function getCartSnapshot() {
    return groupCartItems(cart);
}

document.addEventListener("DOMContentLoaded", () => {
    updateCartCount();
    renderCart();
    setupCartButtons();
});

window.increaseQuantity = increaseQuantity;
window.decreaseQuantity = decreaseQuantity;
window.setupCartButtons = setupCartButtons;
window.addToCart         = addToCart;
window.clearCart         = clearCart;
window.getCartSnapshot   = getCartSnapshot;
window.formatNaira       = formatNaira;