// DOM Elements
const pizzaContainer = document.getElementById("pizzaContainer");
const searchDropdown = document.getElementById("menu");
const searchInput = document.getElementById("search");



function createPizzaCard(pizza) {
    const defaultImage = "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=600&auto=format&fit=crop";
    const images = pizza.images?.length ? pizza.images : [{ image_url: defaultImage }];
    const mainImageUrl = images[0].image_url;

    const card = document.createElement("div");
    card.className = "bg-white rounded-2xl shadow-sm hover:shadow-xl transition-shadow duration-300 border border-gray-100 group flex flex-col overflow-hidden";


    const thumbnailsHTML = images.length > 1 ? `
        <div class="flex gap-2 p-3 overflow-x-auto custom-scrollbar border-b border-gray-50 bg-gray-50/50">
            ${images.map(img => `
                <img
                    src="${img.image_url}"
                    class="thumbnail-btn w-12 h-12 rounded-lg object-cover cursor-pointer border-2 border-transparent hover:border-orange-500 transition-colors"
                    data-image="${img.image_url}"
                    alt="thumbnail"
                >
            `).join("")}
        </div>
    ` : '';

    card.innerHTML = `
        <div class="relative h-56 overflow-hidden bg-gray-100">
            <img
                class="main-image w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                src="${mainImageUrl}"
                alt="${pizza.name}"
            >
        </div>
        
        ${thumbnailsHTML}

        <div class="p-5 flex flex-col flex-grow">
            <h3 class="text-xl font-bold text-gray-900 mb-1">${pizza.name}</h3>
            <p class="text-gray-500 text-sm mb-4 line-clamp-2 flex-grow">${pizza.description || "A delicious freshly baked pizza."}</p>
            
            <div class="mt-auto flex justify-between items-center pt-4 border-t border-gray-50">
                <span class="text-2xl font-black text-orange-600">
                    ₦${Number(pizza.price).toLocaleString()}
                </span>
                <button
                    class="add-to-cart bg-gray-900 hover:bg-orange-600 text-white font-semibold px-4 py-2.5 rounded-xl transition-colors shadow-md active:scale-95 flex items-center gap-2"
                    data-id="${pizza.id}"
                    data-name="${pizza.name}"
                    data-price="${pizza.price}"
                    data-image="${mainImageUrl}"
                >
                Add
                </button>
            </div>
        </div>
    `;

    // Modern, synchronous event binding (No more setTimeout hack!)
    if (images.length > 1) {
        const mainImgEl = card.querySelector('.main-image');
        const thumbs = card.querySelectorAll('.thumbnail-btn');
        
        thumbs.forEach(thumb => {
            thumb.addEventListener("click", (e) => {
                mainImgEl.src = e.target.dataset.image;
            });
        });
    }

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

    searchDropdown.innerHTML = pizzas.map(pizza => {
        const imageUrl = pizza.images?.[0]?.image_url || "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=150&auto=format&fit=crop";
        
        return `
            <div class="flex items-center gap-4 p-4 hover:bg-orange-50 cursor-pointer transition-colors" onclick="window.location.href='pizza.html?id=${pizza.id}'">
                <img src="${imageUrl}" alt="${pizza.name}" class="w-16 h-16 rounded-xl object-cover shadow-sm">
                <div class="flex-grow">
                    <h4 class="text-base font-bold text-gray-900">${pizza.name}</h4>
                    <p class="text-sm text-gray-500 line-clamp-1">${pizza.description}</p>
                </div>
                <span class="font-black text-orange-600 whitespace-nowrap">₦${Number(pizza.price).toLocaleString()}</span>
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
                const response = await fetch(`http://localhost:8080/api/admin/search?q=${encodeURIComponent(query)}`);
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
        // Assuming getMenu() is defined elsewhere in your codebase (e.g., api.js)
        const pizzas = await getMenu(); 

        if (!pizzaContainer) return;
        pizzaContainer.innerHTML = "";

        pizzas.forEach(pizza => {
            pizzaContainer.appendChild(createPizzaCard(pizza));
        });


        if (typeof loadFeaturedPizza === "function") {
            loadFeaturedPizza(pizzas);
        }

    
        if (typeof setupCartButtons === "function") {
            setupCartButtons();
        }

    } catch (error) {
        console.error("Failed to load menu:", error);
        if (pizzaContainer) {
            pizzaContainer.innerHTML = `<div class="col-span-full text-center text-red-500 py-10">Failed to load menu. Please try refreshing.</div>`;
        }
    }
}


document.addEventListener("DOMContentLoaded", loadMenu);