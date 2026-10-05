const API_URL = "";

let currentPizza = null;
let currentImages = [];
let selectedQty = 1;

let selectedSize = "medium";

function currentPriceForSize() {
    if (!currentPizza) return 0;
    switch (selectedSize) {
        case "small": return currentPizza.price_small;
        case "large": return currentPizza.price_large;
        default:      return currentPizza.price_medium;
    }
}

function updatePriceDisplay() {
    if (priceEl) priceEl.textContent = `₦${Number(currentPriceForSize()).toLocaleString("en-NG")}`;
}

function renderSizeSelector() {
    if (!priceEl || !currentPizza) return;

    // Remove any previous selector before re-inserting, in case this ever
    // runs more than once for the same page load.
    document.getElementById("sizeSelector")?.remove();

    const sizes = [
        { key: "small",  label: "Small",  price: currentPizza.price_small },
        { key: "medium", label: "Medium", price: currentPizza.price_medium },
        { key: "large",  label: "Large",  price: currentPizza.price_large },
    ];

    const html = `
        <div id="sizeSelector" class="flex gap-2 my-3">
            ${sizes.map(s => `
                <button
                    type="button"
                    data-size="${s.key}"
                    onclick="selectSize('${s.key}')"
                    class="size-btn px-4 py-2 rounded-xl border-2 text-sm font-semibold transition-all ${
                        s.key === selectedSize
                            ? "border-orange-500 bg-orange-50 text-orange-600"
                            : "border-gray-200 text-gray-600 hover:border-orange-300"
                    }"
                >${s.label}</button>
            `).join("")}
        </div>
    `;

    priceEl.insertAdjacentHTML("afterend", html);
}

function selectSize(size) {
    selectedSize = size;
    updatePriceDisplay();

    document.querySelectorAll(".size-btn").forEach(btn => {
        const isActive = btn.dataset.size === size;
        btn.classList.toggle("border-orange-500", isActive);
        btn.classList.toggle("bg-orange-50", isActive);
        btn.classList.toggle("text-orange-600", isActive);
        btn.classList.toggle("border-gray-200", !isActive);
        btn.classList.toggle("text-gray-600", !isActive);
    });
}
window.selectSize = selectSize;




const mainImageEl = document.getElementById("mainImage");
const thumbsEl    = document.getElementById("thumbs");
const nameEl      = document.getElementById("pizzaName");
const descEl      = document.getElementById("pizzaDesc");
const priceEl     = document.getElementById("pizzaPrice");
const qtyEl       = document.getElementById("qtyDisplay");
const counterEl   = document.getElementById("imgCounter");

function setActiveImage(url, index) {
    if (!mainImageEl) return;
    mainImageEl.src = url;


    document.querySelectorAll(".thumb").forEach((t, i) => {
        t.classList.toggle("active", i === index);
        t.classList.toggle("border-orange-500", i === index);
        t.classList.toggle("border-transparent", i !== index);
    });

 
    if (counterEl && currentImages.length > 1) {
        counterEl.textContent = `${index + 1} / ${currentImages.length}`;
    }
}


window.changeImage = (url, index) => setActiveImage(url, index ?? 0);


function changeQty(delta) {
    selectedQty = Math.max(1, Math.min(10, selectedQty + delta));
    if (qtyEl) qtyEl.textContent = selectedQty;
}
window.changeQty = changeQty;


function addPizzaToCart() {
    if (!currentPizza) return;

    const image = currentImages[0]?.image_url || "";
    const price = currentPriceForSize();

    for (let i = 0; i < selectedQty; i++) {
        addToCart({
            id:    currentPizza.id,
            name:  currentPizza.name,
            size:  selectedSize,
            price,
            image,
        });
    }

   
    const btn = document.getElementById("addToCartBtn");
    if (btn) {
        btn.textContent = `✓ Added ${selectedQty > 1 ? "×" + selectedQty : ""}`;
        btn.classList.add("bg-green-600");
        btn.classList.remove("bg-gray-900");
        setTimeout(() => {
            btn.innerHTML = `
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                Add to Cart
            `;
            btn.classList.remove("bg-green-600");
            btn.classList.add("bg-gray-900");
        }, 1500);
    }

  
    selectedQty = 1;
    if (qtyEl) qtyEl.textContent = "1";
}
window.addPizzaToCart = addPizzaToCart;


async function loadPizzaDetails() {
    const id = new URLSearchParams(window.location.search).get("id");
    if (!id) { showError(); return; }

    try {
        const res = await fetch(`${API_URL}/api/pizzas/${id}`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);

        const data = await res.json();
        if (!data?.pizza) { showError(); return; }

        currentPizza  = data.pizza;
        currentImages = data.images || [];


        if (nameEl)  nameEl.textContent  = currentPizza.name;
        if (descEl)  descEl.textContent  = currentPizza.description || "A freshly baked artisan pizza.";
        renderSizeSelector();
        updatePriceDisplay();

        const fallback = "https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=800&auto=format&fit=crop";

        if (currentImages.length > 0) {
            mainImageEl.src = currentImages[0].image_url;

            if (currentImages.length > 1) {
                counterEl?.classList.remove("hidden");
                counterEl.textContent = `1 / ${currentImages.length}`;

                thumbsEl.innerHTML = currentImages.map((img, i) => `
                    <img
                        src="${img.image_url}"
                        class="thumb w-16 h-16 rounded-xl object-cover border-2 cursor-pointer flex-shrink-0 ${i === 0 ? "border-orange-500" : "border-transparent"}"
                        onclick="changeImage('${img.image_url}', ${i})"
                        alt="View ${i + 1}"
                    >
                `).join("");
            }
        } else {
            mainImageEl.src = fallback;
        }

        document.getElementById("skeletonEl")?.classList.add("hidden");
        document.getElementById("pizzaContent")?.classList.remove("hidden");

    } catch (err) {
        console.error("loadPizzaDetails:", err);
        showError();
    }
}

function showError() {
    document.getElementById("skeletonEl")?.classList.add("hidden");
    document.getElementById("errorEl")?.classList.remove("hidden");
}

document.addEventListener("DOMContentLoaded", loadPizzaDetails);