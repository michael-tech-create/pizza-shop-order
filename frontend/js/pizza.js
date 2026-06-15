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

function changeImage(url) {
    mainImage.src = url;
}

async function loadPizzaDetails() {

    const id = getPizzaId();

    if (!id) {
        alert("No pizza selected!");
        return;
    }

    const res = await fetch(`${API_URL}/api/pizzas/${id}`);
    
    const data = await res.json();

    console.log("API RESPONSE:", data);

    if (!data.pizza) {
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



    thumbs.innerHTML = images.map(img => `
        <img 
            src="${img.image_url}"
            class="w-20 h-20 object-cover rounded-lg border cursor-pointer hover:scale-105 transition"
            onclick="changeImage('${img.image_url}')"
        />
    `).join("");
}




document.addEventListener("DOMContentLoaded", loadPizzaDetails);