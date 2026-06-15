async function getMenu() {
    try {
        const response = await fetch("http://localhost:8080/menu");
        const pizzas = await response.json(); 
        console.log("Parsed Pizza Data:", pizzas); 
        
        return pizzas; 
    } catch (error) {
        console.error("Error loading menu:", error);
        return []; 
    }
}
