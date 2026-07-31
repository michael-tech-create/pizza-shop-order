// =============================================================
// customer-auth.js — Customer account session handling
//
// Include this on any page that needs to know whether a customer is
// logged in (checkout.html, index.html nav, my-orders.html).
// Guest checkout keeps working with zero changes even if this script
// isn't loaded on a page at all — it's purely additive.
//
// Note: uses AUTH_API_URL (not API_URL/API_BASE) deliberately — several
// other files (pizza.js, orders.js) already declare their own top-level
// `const API_URL` / `API_BASE`. Since classic <script> tags share one
// global scope, reusing either name here would throw a duplicate
// declaration error on any page that loads both files together.
// =============================================================

const AUTH_API_URL = "http://localhost:8080"; // TODO: replace with your Render backend URL

const CUSTOMER_TOKEN_KEY  = "mcpizza_customer_token";
const CUSTOMER_NAME_KEY   = "mcpizza_customer_name";
const CUSTOMER_EXPIRY_KEY = "mcpizza_customer_expiry";

function getCustomerToken()  { return localStorage.getItem(CUSTOMER_TOKEN_KEY); }
function getCustomerName()   { return localStorage.getItem(CUSTOMER_NAME_KEY); }
function getCustomerExpiry() { return Number(localStorage.getItem(CUSTOMER_EXPIRY_KEY) || 0); }

function isCustomerLoggedIn() {
    const token = getCustomerToken();
    const expiry = getCustomerExpiry();
    return !!(token && expiry && Date.now() / 1000 < expiry);
}

function saveCustomerSession(data) {
    localStorage.setItem(CUSTOMER_TOKEN_KEY, data.token);
    localStorage.setItem(CUSTOMER_NAME_KEY, data.name);
    localStorage.setItem(CUSTOMER_EXPIRY_KEY, data.expires_at);
}

function clearCustomerSession() {
    localStorage.removeItem(CUSTOMER_TOKEN_KEY);
    localStorage.removeItem(CUSTOMER_NAME_KEY);
    localStorage.removeItem(CUSTOMER_EXPIRY_KEY);
}

function logoutCustomer() {
    clearCustomerSession();
    window.location.href = "index.html";
}

async function customerSignup(name, phone, password) {
    const res = await fetch(`${AUTH_API_URL}/api/customers/signup`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, phone, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Signup failed");
    saveCustomerSession(data);
    return data;
}

async function customerLogin(phone, password) {
    const res = await fetch(`${AUTH_API_URL}/api/customers/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ phone, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Login failed");
    saveCustomerSession(data);
    return data;
}

// Renders a small "Hi, Name" / "Log in" indicator into any element with
// id="customerNavSlot" — pages that don't have that element simply skip
// this, no error.
function renderCustomerNavSlot() {
    const slot = document.getElementById("customerNavSlot");
    if (!slot) return;

    if (isCustomerLoggedIn()) {
        slot.innerHTML = `
            <a href="my-orders.html" class="text-sm font-semibold text-gray-700 hover:text-orange-600">
                Hi, ${getCustomerName()}
            </a>
            <button onclick="logoutCustomer()" class="text-xs text-gray-400 hover:text-red-500 ml-2">Log out</button>
        `;
    } else {
        slot.innerHTML = `
            <a href="userlogin.html" class="text-sm font-semibold text-gray-700 hover:text-orange-600">
                Log in
            </a>
        `;
    }
}

document.addEventListener("DOMContentLoaded", renderCustomerNavSlot);

window.getCustomerToken   = getCustomerToken;
window.getCustomerName    = getCustomerName;
window.isCustomerLoggedIn = isCustomerLoggedIn;
window.customerSignup     = customerSignup;
window.customerLogin      = customerLogin;
window.logoutCustomer     = logoutCustomer;