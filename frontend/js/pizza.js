
const API_URL = "http://localhost:8080";

let currentPizza = null;
let currentImages = [];
let selectedQty = 1;

let selectedSize = "Medium";




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

    for (let i = 0; i < selectedQty; i++) {
        addToCart({
            id:    currentPizza.id,
            name:  currentPizza.name,
            price: currentPizza.price,
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
        if (priceEl) priceEl.textContent = `₦${Number(currentPizza.price).toLocaleString("en-NG")}`;

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