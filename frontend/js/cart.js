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
            <div class="border-b py-2 flex justify-between items-center">
                <div>
                    <p class="font-bold">
                        ${item.name} ${item.quantity > 1 ? `<span class="text-gray-500 font-medium">x${item.quantity}</span>` : ""}
                    </p>
                    <p>$${itemTotal.toFixed(2)}</p>
                </div>
        <div class="flex items-center gap-2">

            <button
             onclick="decreaseQuantity(${item.id})"
            class="bg-green-500 px-2 rounded text-white"
        >
        -
    </button>

            <span>${item.quantity}</span>

         <button
                onclick="increaseQuantity(${item.id})"
              class="bg-green-500 px-2 text-white rounded"
             >
             +
         </button>
        </div>
            </div>
        `;
    }).join("");

    cartItems.innerHTML = cartHTML || `<p class="text-gray-500 text-center py-4 text-sm">Your cart is empty</p>`;
    cartTotal.textContent = total.toFixed(2);
}


// function removeFromCart(id) {

//     const index = cart.findIndex(item => item.id === id);
    
//     if (index !== -1) {
//         cart.splice(index, 1); 
//     }

//     updateCartCount();
//     renderCart();
// }
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