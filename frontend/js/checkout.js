let cart = JSON.parse(localStorage.getItem("cart")) || [];

function renderCheckout() {
    const container = document.getElementById("cartItems");
    const totalEl = document.getElementById("total");

    let total = 0;

    container.innerHTML = cart.map(item => {
        total += item.price * item.quantity;

        return `
            <div class="flex justify-between p-3 border-b">
                <div>
                    <h4>${item.name}</h4>
                    <small>Qty: ${item.quantity}</small>
                </div>
                <div>₦${item.price * item.quantity}</div>
            </div>
        `;
    }).join("");

    totalEl.innerText = total;
}

async function placeOrder() {
    const name = document.getElementById("name").value;
    const address = document.getElementById("address").value;

    if (!name || !address) {
        alert("Please fill all details");
        return;
    }

    const payload = {
        customer_name: name,
        address: address,
        items: cart.map(item => ({
            pizza_id: item.id,
            quantity: item.quantity
        }))
    };

    const res = await fetch("http://localhost:8080/api/orders", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    });

    if (res.ok) {
        alert("Order placed successfully!");
        localStorage.removeItem("cart");
        window.location.href = "/";
    } else {
        alert("Order failed");
    }
}
