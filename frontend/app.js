// ========================================
// CONFIGURATION
// ========================================
const API_URL = "http://localhost:8080";
let token = localStorage.getItem("token");
let currentUser = JSON.parse(localStorage.getItem("user"));

// ========================================
// INITIALISATION
// ========================================
document.addEventListener("DOMContentLoaded", () => {
  checkAuth();
  loadPublicShops();
  setupForms();
});

// ========================================
// AUTHENTIFICATION
// ========================================
function checkAuth() {
  if (token && currentUser) {
    document.getElementById("login-link").style.display = "none";
    document.getElementById("dashboard-link").style.display = "inline";
    document.getElementById("logout-link").style.display = "inline";
  } else {
    document.getElementById("login-link").style.display = "inline";
    document.getElementById("dashboard-link").style.display = "none";
    document.getElementById("logout-link").style.display = "none";
  }
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  token = null;
  currentUser = null;
  checkAuth();
  showSection("home");
}

// ========================================
// NAVIGATION
// ========================================
function showSection(sectionId) {
  document
    .querySelectorAll(".section")
    .forEach((s) => s.classList.remove("active"));
  document.getElementById(sectionId).classList.add("active");

  if (sectionId === "dashboard" && token) {
    loadDashboard();
  }
  if (sectionId === "public") {
    loadPublicShops();
  }
}

function showTab(tabId) {
  document
    .querySelectorAll(".tab")
    .forEach((t) => t.classList.remove("active"));
  document
    .querySelectorAll(".tab-content")
    .forEach((t) => t.classList.remove("active"));

  document
    .querySelector(`[onclick="showTab('${tabId}')"]`)
    .classList.add("active");
  document.getElementById(tabId).classList.add("active");
}

// ========================================
// FORMULAIRES
// ========================================
function setupForms() {
  // Login
  document
    .getElementById("login-form")
    .addEventListener("submit", async (e) => {
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
          showMessage("login-message", "Connexion réussie!", "success");
          setTimeout(() => showSection("dashboard"), 1000);
        } else {
          showMessage("login-message", data.error, "error");
        }
      } catch (err) {
        showMessage("login-message", "Erreur de connexion au serveur", "error");
      }
    });

  // Register
  document
    .getElementById("register-form")
    .addEventListener("submit", async (e) => {
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
          showMessage(
            "register-message",
            "Compte créé avec succès!",
            "success",
          );
          setTimeout(() => showSection("dashboard"), 1000);
        } else {
          showMessage("register-message", data.error, "error");
        }
      } catch (err) {
        showMessage(
          "register-message",
          "Erreur de connexion au serveur",
          "error",
        );
      }
    });

  // Product Form
  document
    .getElementById("product-form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();
      const body = {
        name: document.getElementById("prod-name").value,
        description: document.getElementById("prod-desc").value,
        category: document.getElementById("prod-category").value,
        purchase_price: parseFloat(
          document.getElementById("prod-purchase").value,
        ),
        selling_price: parseFloat(
          document.getElementById("prod-selling").value,
        ),
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
          alert("Produit créé avec succès!");
        } else {
          const data = await res.json();
          alert(data.error);
        }
      } catch (err) {
        alert("Erreur lors de la création");
      }
    });

  // Transaction Form
  document
    .getElementById("transaction-form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();
      const type = document.getElementById("trans-type").value;
      const body = {
        type: type,
        amount: parseFloat(document.getElementById("trans-amount").value),
      };

      if (type === "Sale") {
        body.product_id = parseInt(
          document.getElementById("trans-product").value,
        );
        body.quantity = parseInt(
          document.getElementById("trans-quantity").value,
        );
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
          loadDashboardStats();
          alert("Transaction enregistrée!");
        } else {
          const data = await res.json();
          alert(data.error);
        }
      } catch (err) {
        alert("Erreur lors de la création");
      }
    });

  // User Form
  document.getElementById("user-form").addEventListener("submit", async (e) => {
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
        alert("Utilisateur créé!");
      } else {
        const data = await res.json();
        alert(data.error);
      }
    } catch (err) {
      alert("Erreur lors de la création");
    }
  });
}

// ========================================
// PUBLIC ROUTES
// ========================================
async function loadPublicShops() {
  try {
    const res = await fetch(`${API_URL}/public/shops`);
    const data = await res.json();

    const container = document.getElementById("shops-list");
    if (!data.shops || data.shops.length === 0) {
      container.innerHTML = "<p>Aucun shop disponible</p>";
      return;
    }

    container.innerHTML = data.shops
      .map(
        (shop) => `
            <div class="shop-card" onclick="loadShopProducts(${shop.id}, '${shop.name}')">
                <h3>🏪 ${shop.name}</h3>
                <p>Cliquez pour voir les produits</p>
            </div>
        `,
      )
      .join("");
  } catch (err) {
    console.error("Erreur:", err);
  }
}

async function loadShopProducts(shopId, shopName) {
  try {
    const res = await fetch(`${API_URL}/public/${shopId}/products`);
    const data = await res.json();

    document.getElementById("shops-list").style.display = "none";
    const container = document.getElementById("public-products");
    container.style.display = "grid";

    if (!data.products || data.products.length === 0) {
      container.innerHTML = `
                <button class="btn btn-secondary" onclick="backToShops()">← Retour</button>
                <p>Aucun produit disponible</p>
            `;
      return;
    }

    container.innerHTML = `
            <div style="grid-column: 1 / -1; margin-bottom: 1rem;">
                <button class="btn btn-secondary" onclick="backToShops()">← Retour aux boutiques</button>
                <h3 style="margin-top: 1rem;">Produits de ${shopName}</h3>
            </div>
            ${data.products
              .map(
                (p) => `
                <div class="product-card">
                    <img src="${p.image_url || "https://via.placeholder.com/300x200?text=" + p.name}" alt="${p.name}">
                    <div class="content">
                        <span class="category">${p.category || "Non catégorisé"}</span>
                        <h3>${p.name}</h3>
                        <p>${p.description || ""}</p>
                        <div class="price">${p.selling_price.toLocaleString()} DH</div>
                        <div class="stock ${p.in_stock ? "in-stock" : "out-of-stock"}">
                            ${p.in_stock ? `✅ En stock (${p.stock})` : "❌ Rupture de stock"}
                        </div>
                        ${
                          p.in_stock
                            ? `
                            <a href="${p.whatsapp_link}" target="_blank" class="whatsapp-btn">
                                📱 Commander via WhatsApp
                            </a>
                        `
                            : ""
                        }
                    </div>
                </div>
            `,
              )
              .join("")}
        `;
  } catch (err) {
    console.error("Erreur:", err);
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
  document.getElementById("user-info").textContent =
    `${currentUser.name} (${currentUser.role})`;

  // Masquer l'onglet users pour les Admin
  if (currentUser.role !== "SuperAdmin") {
    document.getElementById("users-tab-btn").style.display = "none";
  } else {
    document.getElementById("users-tab-btn").style.display = "inline";
  }

  await loadDashboardStats();
  await loadProducts();
  await loadTransactions();
  if (currentUser.role === "SuperAdmin") {
    await loadUsers();
  }
}

async function loadDashboardStats() {
  if (currentUser.role !== "SuperAdmin") {
    document.getElementById("stats-cards").innerHTML = `
            <div class="stat-card">
                <div class="label">Rôle</div>
                <div class="value">${currentUser.role}</div>
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

    document.getElementById("stats-cards").innerHTML = `
            <div class="stat-card">
                <div class="label">💰 Ventes Totales</div>
                <div class="value">${d.total_sales.toLocaleString()} DH</div>
            </div>
            <div class="stat-card">
                <div class="label">💸 Dépenses</div>
                <div class="value">${d.total_expenses.toLocaleString()} DH</div>
            </div>
            <div class="stat-card profit">
                <div class="label">📈 Profit Net</div>
                <div class="value">${d.net_profit.toLocaleString()} DH</div>
            </div>
            <div class="stat-card">
                <div class="label">📦 Produits</div>
                <div class="value">${d.total_products}</div>
            </div>
            <div class="stat-card warning">
                <div class="label">⚠️ Stock Faible</div>
                <div class="value">${d.low_stock_products}</div>
            </div>
            <div class="stat-card">
                <div class="label">🔄 Transactions</div>
                <div class="value">${d.transactions.total}</div>
            </div>
        `;
  } catch (err) {
    console.error("Erreur:", err);
  }
}

async function loadProducts() {
  try {
    const res = await fetch(`${API_URL}/products`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    const container = document.getElementById("products-list");
    if (!data.products || data.products.length === 0) {
      container.innerHTML = "<p>Aucun produit</p>";
      return;
    }

    const showPurchasePrice = currentUser.role === "SuperAdmin";

    container.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Nom</th>
                        <th>Catégorie</th>
                        ${showPurchasePrice ? "<th>Prix Achat</th>" : ""}
                        <th>Prix Vente</th>
                        <th>Stock</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.products
                      .map(
                        (p) => `
                        <tr>
                            <td>${p.name}</td>
                            <td>${p.category || "-"}</td>
                            ${showPurchasePrice ? `<td>${p.purchase_price?.toLocaleString() || "-"} DH</td>` : ""}
                            <td>${p.selling_price.toLocaleString()} DH</td>
                            <td class="${p.stock < 5 ? "text-warning" : ""}">${p.stock}</td>
                            <td>
                                <button class="btn btn-danger" onclick="deleteProduct(${p.id})">🗑️</button>
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
  try {
    const res = await fetch(`${API_URL}/transactions`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    const container = document.getElementById("transactions-list");
    if (!data.transactions || data.transactions.length === 0) {
      container.innerHTML = "<p>Aucune transaction</p>";
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
                            <td><span class="badge badge-${t.type.toLowerCase()}">${t.type}</span></td>
                            <td>${t.product?.name || "-"}</td>
                            <td>${t.quantity || "-"}</td>
                            <td>${t.amount.toLocaleString()} DH</td>
                            <td>${new Date(t.created_at).toLocaleDateString()}</td>
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
  try {
    const res = await fetch(`${API_URL}/users`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    const container = document.getElementById("users-list");
    if (!data.users || data.users.length === 0) {
      container.innerHTML = "<p>Aucun utilisateur</p>";
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
                            <td>${u.name}</td>
                            <td>${u.email}</td>
                            <td><span class="badge badge-${u.role.toLowerCase()}">${u.role}</span></td>
                            <td>
                                ${
                                  u.id !== currentUser.id
                                    ? `
                                    <button class="btn btn-danger" onclick="deleteUser(${u.id})">🗑️</button>
                                `
                                    : "(vous)"
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
// ACTIONS
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
    }
  } catch (err) {
    alert("Erreur lors de la suppression");
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
    }
  } catch (err) {
    alert("Erreur lors de la suppression");
  }
}

// ========================================
// MODALS
// ========================================
function showAddProductModal() {
  document.getElementById("product-form").reset();
  document.getElementById("product-modal").classList.add("active");
}

async function showAddTransactionModal() {
  document.getElementById("transaction-form").reset();

  // Charger les produits pour le select
  try {
    const res = await fetch(`${API_URL}/products`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();

    const select = document.getElementById("trans-product");
    select.innerHTML =
      data.products
        ?.map(
          (p) =>
            `<option value="${p.id}">${p.name} (Stock: ${p.stock})</option>`,
        )
        .join("") || "";
  } catch (err) {
    console.error(err);
  }

  document.getElementById("transaction-modal").classList.add("active");
}

function showAddUserModal() {
  document.getElementById("user-form").reset();
  document.getElementById("user-modal").classList.add("active");
}

function closeModal(modalId) {
  document.getElementById(modalId).classList.remove("active");
}

function toggleProductField() {
  const type = document.getElementById("trans-type").value;
  const productField = document.getElementById("product-field");
  const quantityField = document.getElementById("quantity-field");

  if (type === "Sale") {
    productField.style.display = "block";
    quantityField.style.display = "block";
  } else {
    productField.style.display = "none";
    quantityField.style.display = "none";
  }
}

// ========================================
// UTILITAIRES
// ========================================
function showMessage(elementId, message, type) {
  const el = document.getElementById(elementId);
  el.textContent = message;
  el.className = `message ${type}`;
}
