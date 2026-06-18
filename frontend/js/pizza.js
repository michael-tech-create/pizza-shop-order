const API_URL = "http://localhost:8080";

const mainImage = document.getElementById("mainImage");
const thumbs = document.getElementById("thumbs");
const pizzaName = document.getElementById("pizzaName");
const pizzaDesc = document.getElementById("pizzaDesc");
const pizzaPrice = document.getElementById("pizzaPrice");

function getPizzaId() {
    const params = new URLSearchParams(window.location.search);
    return params.get("id");
}

// Fixed: Explicitly attach to window so inline HTML onclick can find it
window.changeImage = function(url) {
    if (mainImage) mainImage.src = url;
};

async function loadPizzaDetails() {
    const id = getPizzaId();

    if (!id) {
        alert("No pizza selected!");
        return;
    }

    try {
        const res = await fetch(`${API_URL}/api/pizzas/${id}`);
        
        if (!res.ok) {
            throw new Error(`Server responded with status: ${res.status}`);
        }

        const data = await res.json();
        console.log("API RESPONSE:", data);

        if (!data || !data.pizza) {
            alert("Pizza not found!");
            return;
        }

        const pizza = data.pizza;
        const images = data.images || [];

        pizzaName.textContent = pizza.name;
        pizzaDesc.textContent = pizza.description;
        pizzaPrice.textContent = `₦${Number(pizza.price).toLocaleString()}`;

        if (images.length > 0) {
            mainImage.src = images[0].image_url;
        } else {
            mainImage.src = "https://images.unsplash.com/photo-1700760934249-93efbb574d23?ixlib=rb-4.1.0&q=85&fm=jpg&crop=entropy&cs=srgb&dl=farhad-ibrahimzade-dFE0FNVd4k0-unsplash.jpg";
        }

        // Fixed: Wrapped the entire map loop inside a single grid container
        thumbs.innerHTML = `
            <div class="grid grid-cols-3 gap-4">
                ${images.map(img => `
                    <img 
                        src="${img.image_url}"
                        class="w-20 h-20 object-cover rounded-lg border cursor-pointer hover:scale-105 transition"
                        onclick="changeImage('${img.image_url}')"
                        alt="Pizza thumbnail"
                    />
                `).join("")}
            </div>
        `;

    } catch (error) {
        console.error("Error loading pizza details:", error);
        alert("Could not load pizza details. Please try again later.");
    }
}

document.addEventListener("DOMContentLoaded", loadPizzaDetails);