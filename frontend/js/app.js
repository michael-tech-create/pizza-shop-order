async function getMenu() {
    try {
        const response = await fetch("/menu");
        const pizzas = await response.json(); 
        console.log("Parsed Pizza Data:", pizzas); 
        
        return pizzas; 
    } catch (error) {
        console.error("Error loading menu:", error);
        return []; 
    }
}
