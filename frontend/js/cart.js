// ============================================================
// cart.js — Cart state management
// Responsibility: in-memory cart array, UI rendering, quantity
// controls. Does NOT talk to the server — that lives in orders.js
// ============================================================

let cart = [];

// ── DOM refs ──
const cartItemsEl  = document.getElementById("cartItems");
const cartTotalEl  = document.getElementById("cartTotal");
const cartCountEl  = document.getElementById("cartCount");

// ── Naira formatter ──
function formatNaira(amount) {
    return "₦" + Number(amount).toLocaleString("en-NG");
}

// ── Add item from any "add-to-cart" button on the page ──
function setupCartButtons() {
    document.querySelectorAll(".add-to-cart").forEach(button => {
        button.addEventListener("click", () => {
            addToCart({
                id:    Number(button.dataset.id),
                name:  button.dataset.name,
                price: Number(button.dataset.price),
                image: button.dataset.image || "",
            });
        });
    });
}

function addToCart(item) {
    cart.push(item);
    updateCartCount();
    renderCart();
}

// ── Group duplicate items so we can show qty badges ──
function groupCartItems(items) {
    const grouped = new Map();
    for (const item of items) {
        const existing = grouped.get(item.id);
        if (existing) {
            existing.quantity++;
        } else {
            grouped.set(item.id, { ...item, quantity: 1 });
        }
    }
    return Array.from(grouped.values());
}

function updateCartCount() {
    if (cartCountEl) cartCountEl.textContent = cart.length;
}

function renderCart() {
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
                <span class="text-4xl mb-2 opacity-60">🍕</span>
                <p class="text-gray-400 font-medium text-sm">Your cart is empty.</p>
            </div>
        `;
    }

    if (cartTotalEl) cartTotalEl.textContent = total.toLocaleString("en-NG");
}

function increaseQuantity(id) {
    const item = cart.find(i => i.id === id);
    if (item) cart.push({ id: item.id, name: item.name, price: item.price, image: item.image });
    updateCartCount();
    renderCart();
}

function decreaseQuantity(id) {
    const index = cart.findIndex(i => i.id === id);
    if (index !== -1) cart.splice(index, 1);
    updateCartCount();
    renderCart();
}

function clearCart() {
    cart = [];
    updateCartCount();
    renderCart();
}

function getCartSnapshot() {
    return groupCartItems(cart);
}

// Expose to window for inline onclick handlers and orders.js
window.increaseQuantity = increaseQuantity;
window.decreaseQuantity = decreaseQuantity;
window.setupCartButtons = setupCartButtons;
window.addToCart = addToCart;
window.clearCart = clearCart;
window.getCartSnapshot = getCartSnapshot;
window.formatNaira = formatNaira;