const pizzaContainer = document.getElementById("pizzaContainer");
const menu = document.getElementById("menu");
const searchInput = document.getElementById("search");


function createPizzaCard(pizza) {
    const card = document.createElement("div");

    card.className =
        "bg-white rounded-2xl shadow-sm border border-amber-200 overflow-hidden flex flex-col transition-all hover:-translate-y-1 hover:shadow-lg hover:shadow-orange-200/50";

    card.onclick = () => {
        window.location.href = `pizza.html?id=${pizza.id}`;
    };

    card.innerHTML = `
        <div class="h-40 bg-gradient-to-br from-orange-200 via-amber-100 to-orange-50 flex items-center justify-center">
            <img
    src="http://localhost:8080${pizza.image_url}"
    alt="${pizza.name}"
    class="w-full h-full object-cover"
/>
        </div>

        <div class="p-5 flex flex-col flex-1">

            <h3 class="text-lg font-bold text-slate-800 mb-1">
                ${pizza.name || "Untitled Pizza"}
            </h3>

            <p class="text-sm text-slate-500 flex-1 mb-4">
                ${pizza.description || "No description available"}
            </p>

            <div class="flex justify-between items-center">

                <span class="text-xl font-bold text-orange-700">
                    ₦${Number(pizza.price || 0).toLocaleString()}
                </span>

                <button
                    class="add-to-cart bg-orange-700 hover:bg-orange-800 text-white px-4 py-2 rounded-xl font-medium"
                    data-id="${pizza.id}"
                    data-name="${pizza.name}"
                    data-price="${pizza.price}"
                >
                    Add To Cart
                </button>

            </div>
        </div>
    `;

    return card;
}

async function loadMenu() {
    try {
        const pizzas = await getMenu();

        pizzaContainer.innerHTML = "";

        pizzas.forEach(pizza => {
            pizzaContainer.appendChild(createPizzaCard(pizza));
        });

        setupCartButtons();

    } catch (error) {
        console.error("Failed to load menu:", error);
    }
}

function renderPizza(pizzas) {
    if (!menu) return;

    if (!pizzas || pizzas.length === 0) {
        menu.innerHTML = `
            <div class="col-span-full flex flex-col items-center justify-center py-16 text-slate-400">
                <span class="text-5xl mb-3">🍕</span>
                <span class="font-medium">No pizzas found.</span>
            </div>
        `;
        return;
    }

    menu.innerHTML = pizzas.map(pizza => `
        <div class="bg-white rounded-2xl shadow-sm border border-amber-200 p-5">
            <h3 class="text-lg font-bold text-slate-800 mb-1">
                ${pizza.name || "Unknown"}
            </h3>

            <p class="text-sm text-slate-500 mb-3">
                ${pizza.description || ""}
            </p>

            <p class="text-xl font-bold text-orange-700">
                ₦${Number(pizza.price || 0).toLocaleString()}
            </p>
        </div>
    `).join("");
}

/* -----------------------------
   SEARCH HANDLER
------------------------------ */
searchInput.addEventListener("input", async (e) => {
    const search = e.target.value;

    try {
        const response = await fetch(
            `http://localhost:8080/api/admin/search?q=${search}`
        );

        const pizzas = await response.json();

        renderPizza(pizzas || []);

    } catch (error) {
        console.error("Search failed:", error);
        renderPizza([]);
    }
});


loadMenu();