let cart = [];

const cartItems = document.getElementById("cartItems");
const cartTotal = document.getElementById("cartTotal");
const cartBtn = document.getElementById("cartBtn");
const cartDrawer = document.getElementById("cart");

cartBtn.addEventListener("click", () => {
    cartDrawer.classList.toggle("hidden");
});

function setupCartButtons() {
    const buttons = document.querySelectorAll(".add-to-cart");

    buttons.forEach(button => {
        button.addEventListener("click", () => {
            const item = {
                id: Number(button.dataset.id),
                name: button.dataset.name,
                price: Number(button.dataset.price)
            };

            cart.push(item);

            updateCartCount();
            renderCart();
        });
    });
}

function updateCartCount() {
    document.getElementById("cartCount").textContent = cart.length;
}

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

function renderCart() {
    const groupedItems = groupCartItems(cart);
    let total = 0;

    const cartHTML = groupedItems.map(item => {
        const itemTotal = item.price * item.quantity;
        total += itemTotal;

        return `
            <div class="py-4 flex justify-between items-center border-b border-gray-100 last:border-b-0 hover:bg-gray-50/50 px-2 rounded-xl transition-colors">
                <div>
                    <p class="font-bold text-gray-800 text-sm flex items-center gap-1.5">
                        ${item.name} 
                        ${item.quantity > 1 ? `<span class="bg-orange-50 text-orange-600 font-semibold text-xs px-2 py-0.5 rounded-md">x${item.quantity}</span>` : ""}
                    </p>
                    <p class="text-sm font-black text-gray-900 mt-0.5">$${itemTotal.toFixed(2)}</p>
                </div>
                
                <!-- Improved Control Buttons Container -->
                <div class="flex items-center gap-2.5 bg-gray-100/80 p-1 rounded-xl border border-gray-200/40">
                    <button
                        onclick="decreaseQuantity(${item.id})"
                        class="bg-white hover:bg-orange-50 hover:text-orange-600 text-gray-600 font-bold w-7 h-7 rounded-lg shadow-sm flex items-center justify-center active:scale-90 transition-all text-xs"
                    >
                        —
                    </button>

                    <span class="text-xs font-bold text-gray-800 px-1 min-w-[12px] text-center">${item.quantity}</span>

                    <button
                        onclick="increaseQuantity(${item.id})"
                        class="bg-white hover:bg-orange-50 hover:text-orange-600 text-gray-600 font-bold w-7 h-7 rounded-lg shadow-sm flex items-center justify-center active:scale-90 transition-all text-sm"
                    >
                        +
                    </button>
                </div>
            </div>
        `;
    }).join("");

    cartItems.innerHTML = cartHTML || `
        <div class="text-center py-12 flex flex-col items-center justify-center">
            <span class="text-4xl mb-2 opacity-60">🍕</span>
            <p class="text-gray-400 font-medium text-sm">Your cart is empty.</p>
        </div>
    `;
    cartTotal.textContent = total.toFixed(2);
}

function increaseQuantity(id) {
    const item = cart.find(item => item.id === id);

    if (item) {
        cart.push({
            id: item.id,
            name: item.name,
            price: item.price
        });
    }

    updateCartCount();
    renderCart();
}

function decreaseQuantity(id) {
    const index = cart.findIndex(item => item.id === id);
    console.log(index)

    if (index !== -1) {
        cart.splice(index, 1)
    }

    updateCartCount();
    renderCart();
}

async function checkout() {
    if (cart.length == 0) {
        alert("cart is empty")
        return;
    }

    const grouped = groupCartItems(cart);
    const firstItem = grouped[0];

    const response = await fetch(
        "http://localhost:8080/api/orders",
        {
            method: "POST", 
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                pizza_id: firstItem.id,
                quantity: firstItem.quantity
            })
        }
    );

    const data = await response.json();
    console.log(data);

    if (response.ok) {
        alert(
            `Order Created!\nTotal: $${data.total_cost}`
        );

        cart = [];
        updateCartCount();
        renderCart();
    } else {
        alert(data.error)
    }
}

window.increaseQuantity = increaseQuantity;
window.decreaseQuantity = decreaseQuantity;
window.checkout = checkout;
