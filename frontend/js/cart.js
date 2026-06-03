let cart = [];

const cartItems = document.getElementById("cartItems")
const cartTotal = document.getElementById('cartTotal')
const cartBtn = document.getElementById("cartBtn");
const cartDrawer = document.getElementById("cart");

cartBtn.addEventListener("click", ()=>{
    cartDrawer.classList.toggle("hidden")
})

function setupCartButtons() {

    const buttons =
        document.querySelectorAll(".add-to-cart");

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

            console.log(cart);
        });

    });

}

function updateCartCount() {

    document.getElementById("cartCount")
        .textContent = cart.length;

}

function renderCart() {
    cartItems.innerHTML = "";

    let total = 0;

    cart.forEach((item, index) => {
        total += item.price;

        cartItems.innerHTML += `
                    <div class="border-b py-2 flex justify-between">
                <div>
                    <p class="font-bold">${item.name}</p>
                    <p>$${item.price}</p>
                </div>

                <button
                    onclick="removeFromCart(${index})"
                    class="text-red-500">
                    x
                </button>
            </div>
        `;

    });

    cartTotal.textContent = total;
}

function removeFromCart(index) {
    cart.splice(index, 1)

    updateCartCount();
    renderCart()
}