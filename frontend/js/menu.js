const pizzaContainer =
    document.getElementById("pizzaContainer");

async function loadMenu() {

    const pizzas = await getMenu();

    pizzaContainer.innerHTML = "";

    pizzas.forEach(pizza => {

        const card = document.createElement("div");

        card.className =
            "bg-white rounded-lg shadow p-4";

        card.innerHTML = `
            <h3 class="text-xl font-bold mb-2">
                ${pizza.name}
            </h3>

            <p class="text-gray-600 mb-4">
                ${pizza.description}
            </p>

            <div class="flex justify-between items-center">

                <span class="font-bold">
                    $${pizza.price}
                </span>

                <button
                    class="add-to-cart bg-orange-700 text-white px-3 py-2 rounded"
                    data-id="${pizza.id}"
                    data-name="${pizza.name}"
                    data-price="${pizza.price}"
                >
                    Add To Cart
                </button>

            </div>
        `;

        pizzaContainer.appendChild(card);
    });

    setupCartButtons();
}

loadMenu();