// DOM Elements
const pizzaContainer = document.getElementById("pizzaContainer");
const searchDropdown = document.getElementById("menu");
const searchInput = document.getElementById("search");
const categoryFiltersEl = document.getElementById("categoryFilters");

let allPizzas = [];
let activeCategoryId = null; // null = "All"

// ── Category filter tabs ──
async function loadCategoryFilters() {
    if (!categoryFiltersEl) return;
    try {
        const res = await fetch("/api/categories");
        if (!res.ok) return;
        const categories = await res.json();
        renderCategoryFilters(categories);
    } catch (e) {
        console.error("Failed to load categories:", e);
    }
}

function renderCategoryFilters(categories) {
    const tabClass = (isActive) => [
        "px-4 py-2 rounded-full text-sm font-semibold transition-all flex items-center gap-2",
        isActive
            ? "bg-orange-600 text-white shadow-md shadow-orange-600/20"
            : "bg-white text-gray-600 border border-gray-200 hover:border-orange-300 hover:text-orange-600"
    ].join(" ");

    const allBtn = `<button data-cat="" class="cat-filter-btn ${tabClass(activeCategoryId === null)}">All</button>`;
    const catBtns = categories.map(c => `
        <button data-cat="${c.id}" class="cat-filter-btn ${tabClass(activeCategoryId === c.id)}">
            ${c.image_url ? `<img src="${c.image_url}" class="w-5 h-5 rounded-full object-cover" onerror="this.remove()">` : ""}
            ${c.name}
        </button>
    `).join("");

    categoryFiltersEl.innerHTML = allBtn + catBtns;

    categoryFiltersEl.querySelectorAll(".cat-filter-btn").forEach(btn => {
        btn.addEventListener("click", () => {
            const val = btn.dataset.cat;
            activeCategoryId = val === "" ? null : Number(val);
            renderCategoryFilters(categories); // refresh active styling
            applyMenuFilter();
        });
    });
}

function applyMenuFilter() {
    if (!pizzaContainer) return;
    const filtered = activeCategoryId === null
        ? allPizzas
        : allPizzas.filter(p => p.category_id === activeCategoryId);

    pizzaContainer.innerHTML = "";
    filtered.forEach(pizza => pizzaContainer.appendChild(createPizzaCard(pizza)));

    if (typeof setupCartButtons === "function") setupCartButtons();
}

function createPizzaCard(pizza) {
    const defaultImage = "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=600&auto=format&fit=crop";
    const images = pizza.images?.length ? pizza.images : [{ image_url: defaultImage }];
    const mainImageUrl = images[0].image_url;

    const card = document.createElement("div");
    card.className = "bg-white rounded-xl sm:rounded-2xl shadow-sm hover:shadow-xl transition-shadow duration-300 border border-gray-100 group flex flex-col overflow-hidden cursor-pointer";
    card.dataset.id = pizza.id;

    const thumbnailsHTML = images.length > 1 ? `
        <div class="flex gap-2 p-2 sm:p-3 overflow-x-auto custom-scrollbar border-b border-gray-50 bg-gray-50/50">
            ${images.map(img => `
                <img
                    src="${img.image_url}"
                    class="thumbnail-btn w-8 h-8 sm:w-12 sm:h-12 rounded-md sm:rounded-lg object-cover cursor-pointer border-2 border-transparent hover:border-orange-500 transition-colors flex-shrink-0"
                    data-image="${img.image_url}"
                    alt="thumbnail"
                >
            `).join("")}
        </div>
    ` : '';

    card.innerHTML = `
        <div class="relative h-32 sm:h-56 overflow-hidden bg-gray-100">
            <img
                class="main-image w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                src="${mainImageUrl}"
                alt="${pizza.name}"
            >
        </div>
        
        ${thumbnailsHTML}

        <div class="p-3 sm:p-5 flex flex-col flex-grow">
            <h3 class="text-base sm:text-xl font-bold text-gray-900 mb-1 truncate">${pizza.name}</h3>
            <p class="text-gray-500 text-xs sm:text-sm mb-3 sm:mb-4 line-clamp-2 flex-grow">${pizza.description || "A delicious freshly baked pizza."}</p>
            
            <div class="mt-auto flex flex-col xl:flex-row justify-between items-start xl:items-center pt-3 sm:pt-4 border-t border-gray-50 gap-2 xl:gap-0">
                <span class="text-lg sm:text-2xl font-black text-orange-600">
                    From ₦${Number(pizza.price_small).toLocaleString()}
                </span>
                <button
                    class="add-to-cart bg-gray-900 hover:bg-orange-600 text-white font-semibold px-3 py-2 sm:px-4 sm:py-2.5 rounded-lg sm:rounded-xl transition-colors shadow-md active:scale-95 flex items-center justify-center gap-2 w-full xl:w-auto text-sm sm:text-base"
                    data-id="${pizza.id}"
                    data-name="${pizza.name}"
                    data-size="medium"
                    data-price="${pizza.price_medium}"
                    data-image="${mainImageUrl}"
                    title="Adds a Medium — open the pizza for other sizes"
                >
                Add
                </button>
            </div>
        </div>
    `;

    // Modern, synchronous event binding
    if (images.length > 1) {
        const mainImgEl = card.querySelector('.main-image');
        const thumbs = card.querySelectorAll('.thumbnail-btn');
        
        thumbs.forEach(thumb => {
            thumb.addEventListener("click", (e) => {
                e.stopPropagation();
                mainImgEl.src = e.target.dataset.image;
            });
        });
    }

    // Clicking anywhere on the card opens the pizza detail page
    card.addEventListener("click", (e) => {
        if (e.target.closest(".add-to-cart") || e.target.closest(".thumbnail-btn")) return;
        window.location.href = `pizza.html?id=${pizza.id}`;
    });

    return card;
}

// --- 2. SEARCH DROPDOWN RENDERER ---
function renderSearchDropdown(pizzas) {
    if (!searchDropdown) return;

    if (!pizzas || pizzas.length === 0) {
        searchDropdown.innerHTML = `
            <div class="p-8 text-center text-gray-400">
                <span class="text-3xl block mb-2">0</span>
                <span class="font-medium text-sm">No pizzas found matching your search.</span>
            </div>
        `;
        return;
    }

    searchDropdown.className = "absolute top-full left-0 z-50 w-full bg-white rounded-xl shadow-lg border border-gray-100 overflow-hidden mt-2";

    searchDropdown.innerHTML = pizzas.map(pizza => {
        const imageUrl = pizza.images?.[0]?.image_url || "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=150&auto=format&fit=crop";
        
        return `
            <div class="flex items-center gap-4 p-4 hover:bg-orange-50 border-b border-gray-50 last:border-none cursor-pointer transition-colors" onclick="window.location.href='pizza.html?id=${pizza.id}'">
                <img src="${imageUrl}" alt="${pizza.name}" class="w-16 h-16 rounded-xl object-cover shadow-sm">
                <div class="flex-grow min-w-0">
                    <h4 class="text-base font-bold text-gray-900 truncate">${pizza.name}</h4>
                    <p class="text-sm text-gray-500 line-clamp-1">${pizza.description}</p>
                </div>
                <span class="font-black text-orange-600 whitespace-nowrap">From ₦${Number(pizza.price_small).toLocaleString()}</span>
            </div>
        `;
    }).join("");
}

let searchTimeout;
if (searchInput) {
    searchInput.addEventListener("input", (e) => {
        const query = e.target.value.trim();
        
        clearTimeout(searchTimeout); 

        if (query.length === 0) {
            searchDropdown.innerHTML = ""; // Hide dropdown if empty
            return;
        }

        searchTimeout = setTimeout(async () => {
            try {
                const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                
                if (!response.ok) throw new Error("Search failed");
                const pizzas = await response.json();
                renderSearchDropdown(pizzas);
            } catch (error) {
                console.error("Search Error:", error);
                renderSearchDropdown([]);
            }
        }, 300); // 300ms debounce
    });

    // Close dropdown when clicking outside
    document.addEventListener("click", (e) => {
        if (!searchInput.contains(e.target) && !searchDropdown.contains(e.target)) {
            searchDropdown.innerHTML = ""; 
        }
    });
}

// Main Load Function
async function loadMenu() {
    try {
        const pizzas = await getMenu(); 

        if (!pizzaContainer) return;
        
        // 1. Ensure the container has the responsive grid classes!
        // Grid 2 columns by default (mobile), 3 on medium, 4 on large screens.
        pizzaContainer.className = "grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-6";

        allPizzas = pizzas;
        applyMenuFilter(); // renders using the current category filter (default: All)

        loadCategoryFilters();

        if (typeof loadFeaturedPizza === "function") {
            loadFeaturedPizza(pizzas);
        }
    
        if (typeof setupCartButtons === "function") {
            setupCartButtons();
        }

    } catch (error) {
        console.error("Failed to load menu:", error);
        if (pizzaContainer) {
            // Restore default view on error so the message isn't constrained to a tiny grid cell
            pizzaContainer.className = "block";
            pizzaContainer.innerHTML = `<div class="text-center text-red-500 py-10">Failed to load menu. Please try refreshing.</div>`;
        }
    }
}

document.addEventListener("DOMContentLoaded", loadMenu);