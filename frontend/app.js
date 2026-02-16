// ========================================
// CONFIGURATION
// ========================================
const API_URL = "http://localhost:8080";
let token = localStorage.getItem("token");
let currentUser = JSON.parse(localStorage.getItem("user"));

// ========================================
// INITIALIZATION
// ========================================
document.addEventListener("DOMContentLoaded", () => {
  initApp();
});

function initApp() {
  checkAuth();
  setupEventListeners();
  loadPublicShops();
  setActiveNavLink("home");
}

// ========================================
// AUTHENTICATION
// ========================================
function checkAuth() {
  const loginLink = document.getElementById("login-link");
  const dashboardLink = document.getElementById("dashboard-link");
  const logoutLink = document.getElementById("logout-link");
  const userCard = document.getElementById("user-card");

  if (token && currentUser) {
    loginLink.style.display = "none";
    dashboardLink.style.display = "flex";
    logoutLink.style.display = "flex";
    userCard.style.display = "flex";

    document.getElementById("user-avatar").textContent = currentUser.name
      .charAt(0)
      .toUpperCase();
    document.getElementById("user-name").textContent = currentUser.name;
    document.getElementById("user-role").textContent = currentUser.role;
  } else {
    loginLink.style.display = "flex";
    dashboardLink.style.display = "none";
    logoutLink.style.display = "none";
    userCard.style.display = "none";
  }
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  token = null;
  currentUser = null;
  checkAuth();
  showSection("home");
  showToast("Déconnexion réussie", "success");
}

// ========================================
// NAVIGATION
// ========================================
function showSection(sectionId) {
  document
    .querySelectorAll(".section")
    .forEach((s) => s.classList.remove("active"));
  document.getElementById(sectionId).classList.add("active");
  setActiveNavLink(sectionId);

  if (sectionId === "dashboard" && token) {
    loadDashboard();
  }
  if (sectionId === "public") {
    loadPublicShops();
  }

  // Close sidebar on mobile
  document.getElementById("sidebar").classList.remove("open");
}

function setActiveNavLink(sectionId) {
  document
    .querySelectorAll(".nav-link")
    .forEach((link) => link.classList.remove("active"));
  const activeLink = document.querySelector(
    `.nav-link[onclick="showSection('${sectionId}')"]`,
  );
  if (activeLink) activeLink.classList.add("active");
}

function toggleSidebar() {
  document.getElementById("sidebar").classList.toggle("open");
}

// ========================================
// EVENT LISTENERS
// ========================================
function setupEventListeners() {
  // Login form
  document.getElementById("login-form").addEventListener("submit", handleLogin);

  // Register form
  document
    .getElementById("register-form")
    .addEventListener("submit", handleRegister);

  // Product form
  document
    .getElementById("product-form")
    .addEventListener("submit", handleProductSubmit);

  // Transaction form
  document
    .getElementById("transaction-form")
    .addEventListener("submit", handleTransactionSubmit);

  // User form
  document
    .getElementById("user-form")
    .addEventListener("submit", handleUserSubmit);

  // Tabs
  document.querySelectorAll(".tab").forEach((tab) => {
    tab.addEventListener("click", () => {
      const tabId = tab.dataset.tab;
      document
        .querySelectorAll(".tab")
        .forEach((t) => t.classList.remove("active"));
      document
        .querySelectorAll(".tab-content")
        .forEach((c) => c.classList.remove("active"));
      tab.classList.add("active");
      document.getElementById(tabId).classList.add("active");
    });
  });
}

// ========================================
// AUTH HANDLERS
// ========================================
async function handleLogin(e) {
  e.preventDefault();

  const email = document.getElementById("login-email").value;
  const password = document.getElementById("login-password").value;

  try {
    const res = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    const data = await res.json();

    if (res.ok) {
      token = data.token;
      currentUser = data.user;
      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(currentUser));
      checkAuth();
      showSection("dashboard");
      showToast("Connexion réussie !", "success");
    } else {
      showMessage("login-message", data.error, "error");
    }
  } catch (err) {
    showMessage("login-message", "Erreur de connexion au serveur", "error");
  }
}

async function handleRegister(e) {
  e.preventDefault();

  const body = {
    name: document.getElementById("reg-name").value,
    email: document.getElementById("reg-email").value,
    password: document.getElementById("reg-password").value,
    role: document.getElementById("reg-role").value,
    shop_name: document.getElementById("reg-shop").value,
    whatsapp_number: document.getElementById("reg-whatsapp").value,
  };

  try {
    const res = await fetch(`${API_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const data = await res.json();

    if (res.ok) {
      token = data.token;
      currentUser = data.user;
      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(currentUser));
      checkAuth();
      showSection("dashboard");
      showToast("Compte créé avec succès !", "success");
    } else {
      showMessage("register-message", data.error, "error");
    }
  } catch (err) {
    showMessage("register-message", "Erreur de connexion au serveur", "error");
  }
}

// ========================================
// PUBLIC ROUTES
// ========================================
async function loadPublicShops() {
  const container = document.getElementById("shops-list");

  try {
    const res = await fetch(`${API_URL}/public/shops`);
    const data = await res.json();

    if (!data.shops || data.shops.length === 0) {
      container.innerHTML = `
                <div class="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                    </svg>
                    <h4>Aucune boutique disponible</h4>
                    <p>Les boutiques apparaîtront ici</p>
                </div>
            `;
      return;
    }

    container.innerHTML = data.shops
      .map(
        (shop) => `
            <div class="shop-card" onclick="loadShopProducts(${shop.id}, '${shop.name}')">
                <div class="shop-icon">🏪</div>
                <h3>${shop.name}</h3>
                <p>Cliquez pour voir les produits</p>
            </div>
        `,
      )
      .join("");
  } catch (err) {
    console.error("Erreur:", err);
    container.innerHTML = `
            <div class="empty-state">
                <h4>Erreur de chargement</h4>
                <p>Impossible de charger les boutiques</p>
            </div>
        `;
  }
}

async function loadShopProducts(shopId, shopName) {
  document.getElementById("shops-list").style.display = "none";
  document.getElementById("public-products").style.display = "block";

  const shopInfo = document.getElementById("shop-info");
  const grid = document.getElementById("products-grid");

  try {
    const res = await fetch(`${API_URL}/public/${shopId}/products`);
    const data = await res.json();

    shopInfo.innerHTML = `
            <h2>${shopName}</h2>
            <p>${data.products?.length || 0} produit(s) disponible(s)</p>
        `;

    if (!data.products || data.products.length === 0) {
      grid.innerHTML = `
                <div class="empty-state" style="grid-column: 1 / -1;">
                    <h4>Aucun produit disponible</h4>
                    <p>Cette boutique n'a pas encore de produits</p>
                </div>
            `;
      return;
    }

    grid.innerHTML = data.products
      .map(
        (p) => `
            <div class="product-card">
                <img class="product-image" src="${getProductImage(p)}" alt="${p.name}" 
                     onerror="this.src='https://ui-avatars.com/api/?name=${encodeURIComponent(p.name)}&background=6366f1&color=fff&size=200'">
                <div class="product-content">
                    <span class="product-category">${p.category || "Non catégorisé"}</span>
                    <h3 class="product-name">${p.name}</h3>
                    <p class="product-description">${p.description || "Aucune description"}</p>
                    <div class="product-price">${formatPrice(p.selling_price)} <span>DH</span></div>
                    <div class="product-stock ${p.in_stock ? "in-stock" : "out-of-stock"}">
                        <span class="stock-dot"></span>
                        ${p.in_stock ? `En stock (${p.stock})` : "Rupture de stock"}
                    </div>
                    ${
                      p.in_stock
                        ? `
                        <a href="${p.whatsapp_link}" target="_blank" class="whatsapp-btn">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
                            </svg>
                            Commander via WhatsApp
                        </a>
                    `
                        : `<button class="btn btn-outline btn-full" disabled>Indisponible</button>`
                    }
                </div>
            </div>
        `,
      )
      .join("");
  } catch (err) {
    console.error("Erreur:", err);
    grid.innerHTML = `
            <div class="empty-state" style="grid-column: 1 / -1;">
                <h4>Erreur de chargement</h4>
                <p>Impossible de charger les produits</p>
            </div>
        `;
  }
}

function backToShops() {
  document.getElementById("shops-list").style.display = "grid";
  document.getElementById("public-products").style.display = "none";
  loadPublicShops();
}

// ========================================
// DASHBOARD
// ========================================
async function loadDashboard() {
  document.getElementById("dashboard-welcome").textContent =
    `Bienvenue, ${currentUser.name}`;
  document.getElementById("dashboard-user-badge").innerHTML = `
        <span class="badge badge-primary">${currentUser.role}</span>
    `;

  // Hide users tab for Admin
  const usersTabBtn = document.getElementById("users-tab-btn");
  if (currentUser.role !== "SuperAdmin") {
    usersTabBtn.style.display = "none";
  } else {
    usersTabBtn.style.display = "flex";
  }

  await loadDashboardStats();
  await loadProducts();
  await loadTransactions();

  if (currentUser.role === "SuperAdmin") {
    await loadUsers();
  }
}

async function loadDashboardStats() {
  const container = document.getElementById("stats-grid");

  if (currentUser.role !== "SuperAdmin") {
    container.innerHTML = `
            <div class="stat-card">
                <div class="stat-icon primary">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                        <circle cx="12" cy="7" r="4"/>
                    </svg>
                </div>
                <div class="stat-value">${currentUser.role}</div>
                <div class="stat-label">Votre rôle</div>
            </div>
        `;
    return;
  }

  try {
    const res = await fetch(`${API_URL}/reports/dashboard`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();
    const d = data.dashboard;

    container.innerHTML = `
            <div class="stat-card">
                <div class="stat-icon primary">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
                    </svg>
                </div>
                <div class="stat-value">${formatPrice(d.total_sales)}</div>
                <div class="stat-label">Ventes totales (DH)</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon danger">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="23,6 13.5,15.5 8.5,10.5 1,18"/><polyline points="17,6 23,6 23,12"/>
                    </svg>
                </div>
                <div class="stat-value">${formatPrice(d.total_expenses)}</div>
                <div class="stat-label">Dépenses (DH)</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon success">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="23,6 13.5,15.5 8.5,10.5 1,18"/><polyline points="17,6 23,6 23,12"/>
                    </svg>
                </div>
                <div class="stat-value">${formatPrice(d.net_profit)}</div>
                <div class="stat-label">Profit net (DH)</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon primary">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                    </svg>
                </div>
                <div class="stat-value">${d.total_products}</div>
                <div class="stat-label">Produits</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon warning">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                        <line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>
                    </svg>
                </div>
                <div class="stat-value">${d.low_stock_products}</div>
                <div class="stat-label">Stock faible</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon primary">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/>
                    </svg>
                </div>
                <div class="stat-value">${d.transactions?.total || 0}</div>
                <div class="stat-label">Transactions</div>
            </div>
        `;
  } catch (err) {
    console.error("Erreur:", err);
  }
}

async function loadProducts() {
  const container = document.getElementById("products-table");

  try {
    const res = await fetch(`${API_URL}/products`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    if (!data.products || data.products.length === 0) {
      container.innerHTML = `
                <div class="empty-state">
                    <h4>Aucun produit</h4>
                    <p>Commencez par ajouter votre premier produit</p>
                </div>
            `;
      return;
    }

    const showPurchasePrice = currentUser.role === "SuperAdmin";

    container.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Produit</th>
                        <th>Catégorie</th>
                        ${showPurchasePrice ? "<th>Prix achat</th>" : ""}
                        <th>Prix vente</th>
                        <th>Stock</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.products
                      .map(
                        (p) => `
                        <tr>
                            <td>
                                <div style="display: flex; align-items: center; gap: 0.75rem;">
                                    <img class="table-image" src="${getProductImage(p)}" alt="${p.name}"
                                         onerror="this.src='https://ui-avatars.com/api/?name=${encodeURIComponent(p.name)}&background=6366f1&color=fff&size=48'">
                                    <div>
                                        <strong>${p.name}</strong>
                                        <div style="font-size: 0.8rem; color: var(--gray-500);">${(p.description || "").substring(0, 30)}...</div>
                                    </div>
                                </div>
                            </td>
                            <td><span class="badge badge-primary">${p.category || "N/A"}</span></td>
                            ${showPurchasePrice ? `<td>${formatPrice(p.purchase_price || 0)} DH</td>` : ""}
                            <td><strong>${formatPrice(p.selling_price)} DH</strong></td>
                            <td>
                                <span class="badge ${p.stock < 5 ? "badge-warning" : "badge-success"}">
                                    ${p.stock} ${p.stock < 5 ? "⚠️" : ""}
                                </span>
                            </td>
                            <td>
                                <button class="btn btn-sm btn-danger" onclick="deleteProduct(${p.id})">Supprimer</button>
                            </td>
                        </tr>
                    `,
                      )
                      .join("")}
                </tbody>
            </table>
        `;
  } catch (err) {
    console.error("Erreur:", err);
  }
}

async function loadTransactions() {
  const container = document.getElementById("transactions-table");

  try {
    const res = await fetch(`${API_URL}/transactions`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    if (!data.transactions || data.transactions.length === 0) {
      container.innerHTML = `
                <div class="empty-state">
                    <h4>Aucune transaction</h4>
                    <p>Les transactions apparaîtront ici</p>
                </div>
            `;
      return;
    }

    container.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Produit</th>
                        <th>Quantité</th>
                        <th>Montant</th>
                        <th>Date</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.transactions
                      .map(
                        (t) => `
                        <tr>
                            <td>
                                <span class="badge ${t.type === "Sale" ? "badge-success" : t.type === "Expense" ? "badge-danger" : "badge-warning"}">
                                    ${t.type === "Sale" ? "Vente" : t.type === "Expense" ? "Dépense" : "Retrait"}
                                </span>
                            </td>
                            <td>${t.product?.name || "-"}</td>
                            <td>${t.quantity || "-"}</td>
                            <td><strong>${formatPrice(t.amount)} DH</strong></td>
                            <td>${formatDate(t.created_at)}</td>
                        </tr>
                    `,
                      )
                      .join("")}
                </tbody>
            </table>
        `;
  } catch (err) {
    console.error("Erreur:", err);
  }
}

async function loadUsers() {
  const container = document.getElementById("users-table");

  try {
    const res = await fetch(`${API_URL}/users`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    if (!data.users || data.users.length === 0) {
      container.innerHTML = `
                <div class="empty-state">
                    <h4>Aucun utilisateur</h4>
                    <p>Ajoutez des utilisateurs à votre équipe</p>
                </div>
            `;
      return;
    }

    container.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Nom</th>
                        <th>Email</th>
                        <th>Rôle</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.users
                      .map(
                        (u) => `
                        <tr>
                            <td><strong>${u.name}</strong></td>
                            <td>${u.email}</td>
                            <td>
                                <span class="badge ${u.role === "SuperAdmin" ? "badge-primary" : "badge-success"}">
                                    ${u.role}
                                </span>
                            </td>
                            <td>
                                ${
                                  u.id !== currentUser.id
                                    ? `<button class="btn btn-sm btn-danger" onclick="deleteUser(${u.id})">Supprimer</button>`
                                    : '<span style="color: var(--gray-400);">Vous</span>'
                                }
                            </td>
                        </tr>
                    `,
                      )
                      .join("")}
                </tbody>
            </table>
        `;
  } catch (err) {
    console.error("Erreur:", err);
  }
}

// ========================================
// FORM HANDLERS
// ========================================
async function handleProductSubmit(e) {
  e.preventDefault();

  const body = {
    name: document.getElementById("prod-name").value,
    description: document.getElementById("prod-desc").value,
    category: document.getElementById("prod-category").value,
    purchase_price: parseFloat(document.getElementById("prod-purchase").value),
    selling_price: parseFloat(document.getElementById("prod-selling").value),
    stock: parseInt(document.getElementById("prod-stock").value),
    image_url: document.getElementById("prod-image").value,
  };

  try {
    const res = await fetch(`${API_URL}/products`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });

    if (res.ok) {
      closeModal("product-modal");
      loadProducts();
      loadDashboardStats();
      showToast("Produit créé avec succès", "success");
      e.target.reset();
    } else {
      const data = await res.json();
      showToast(data.error, "error");
    }
  } catch (err) {
    showToast("Erreur lors de la création", "error");
  }
}

async function handleTransactionSubmit(e) {
  e.preventDefault();

  const type = document.getElementById("trans-type").value;
  const body = {
    type: type,
    amount: parseFloat(document.getElementById("trans-amount").value),
  };

  if (type === "Sale") {
    body.product_id = parseInt(document.getElementById("trans-product").value);
    body.quantity = parseInt(document.getElementById("trans-quantity").value);
  }

  try {
    const res = await fetch(`${API_URL}/transactions`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });

    if (res.ok) {
      closeModal("transaction-modal");
      loadTransactions();
      loadProducts();
      loadDashboardStats();
      showToast("Transaction enregistrée", "success");
      e.target.reset();
    } else {
      const data = await res.json();
      showToast(data.error, "error");
    }
  } catch (err) {
    showToast("Erreur lors de la création", "error");
  }
}

async function handleUserSubmit(e) {
  e.preventDefault();

  const body = {
    name: document.getElementById("user-name").value,
    email: document.getElementById("user-email").value,
    password: document.getElementById("user-password").value,
    role: document.getElementById("user-role").value,
  };

  try {
    const res = await fetch(`${API_URL}/users`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });

    if (res.ok) {
      closeModal("user-modal");
      loadUsers();
      showToast("Utilisateur créé", "success");
      e.target.reset();
    } else {
      const data = await res.json();
      showToast(data.error, "error");
    }
  } catch (err) {
    showToast("Erreur lors de la création", "error");
  }
}

// ========================================
// DELETE ACTIONS
// ========================================
async function deleteProduct(id) {
  if (!confirm("Supprimer ce produit ?")) return;

  try {
    const res = await fetch(`${API_URL}/products/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    });

    if (res.ok) {
      loadProducts();
      loadDashboardStats();
      showToast("Produit supprimé", "success");
    }
  } catch (err) {
    showToast("Erreur lors de la suppression", "error");
  }
}

async function deleteUser(id) {
  if (!confirm("Supprimer cet utilisateur ?")) return;

  try {
    const res = await fetch(`${API_URL}/users/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    });

    if (res.ok) {
      loadUsers();
      showToast("Utilisateur supprimé", "success");
    }
  } catch (err) {
    showToast("Erreur lors de la suppression", "error");
  }
}

// ========================================
// MODALS
// ========================================
function showModal(modalId) {
  document.getElementById(modalId).classList.add("active");

  if (modalId === "transaction-modal") {
    loadProductsForSelect();
  }
}

function closeModal(modalId) {
  document.getElementById(modalId).classList.remove("active");
}

async function loadProductsForSelect() {
  try {
    const res = await fetch(`${API_URL}/products`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    const select = document.getElementById("trans-product");
    select.innerHTML = (data.products || [])
      .map(
        (p) => `<option value="${p.id}">${p.name} (Stock: ${p.stock})</option>`,
      )
      .join("");
  } catch (err) {
    console.error(err);
  }
}

function toggleTransactionFields() {
  const type = document.getElementById("trans-type").value;
  const saleFields = document.getElementById("sale-fields");
  saleFields.style.display = type === "Sale" ? "block" : "none";
}

// ========================================
// UTILITIES
// ========================================
function getProductImage(product) {
  if (product.image_url && product.image_url.startsWith("http")) {
    return product.image_url;
  }
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(product.name)}&background=6366f1&color=fff&size=200&font-size=0.33`;
}

function formatPrice(price) {
  return new Intl.NumberFormat("fr-MA").format(price || 0);
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function showMessage(elementId, message, type) {
  const el = document.getElementById(elementId);
  el.textContent = message;
  el.className = `message ${type}`;
  setTimeout(() => {
    el.className = "message";
  }, 5000);
}

function showToast(message, type = "success") {
  const container = document.getElementById("toast-container");
  const toast = document.createElement("div");
  toast.className = `toast ${type}`;
  toast.innerHTML = `<span class="toast-message">${message}</span>`;
  container.appendChild(toast);

  setTimeout(() => {
    toast.remove();
  }, 4000);
}
