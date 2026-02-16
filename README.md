# 🏪 Electronic Shop Management API

## Bootcamp Go 2024 - Projet Final

API Backend en Go + Frontend pour la gestion multi-boutiques d'électronique avec isolation complète des données, gestion des rôles et intégration WhatsApp.

---

## 👥 Équipe

| Membre | Rôle |
|--------|------|
| **Boudabous Elyes** | Développeur Full Stack |
| **Yasmine Aoudjit** | Développeuse Frontend & Documentation |
| **Mohamed Amine Dhaoui** | Développeur Full-Stack |
| **Mohamed Amine Ourraki** | Développeur Full-Stack |

**Encadrant** : Mr. Mounir Aziz

---

## 🛠️ Stack Technique

- **Backend** : Go 1.22 + Gin Framework + GORM
- **Base de données** : SQLite (pure Go, sans CGO)
- **Authentification** : JWT (JSON Web Tokens)
- **Frontend** : HTML5 + CSS3 + JavaScript Vanilla
- **Containerisation** : Docker + Docker Compose

---

## 📋 Prérequis

- **Docker** et **Docker Compose** installés
- OU **Go 1.22+** installé localement

---

## 🚀 Installation & Exécution

### Option 1 : Avec Docker (Recommandé) ✅
```bash
# 1. Cloner le repository
git clone https://github.com/votre-repo/electronic-shop-api.git

# 2. Entrer dans le dossier
cd electronic-shop-api

# 3. Lancer avec Docker Compose
docker-compose up --build

# 4. Accéder à l'application
# API Backend : http://localhost:8080
# Frontend : http://localhost:3000
```

### Option 2 : Sans Docker
```bash
# 1. Cloner le repository
git clone https://github.com/votre-repo/electronic-shop-api.git

# 2. Entrer dans le dossier
cd electronic-shop-api

# 3. Installer les dépendances Go
go mod download

# 4. Lancer le backend
go run main.go

# 5. Ouvrir le frontend
# Ouvrir frontend/index.html dans un navigateur
# OU utiliser Live Server de VS Code

# 6. Accéder à l'API
# http://localhost:8080
```

---

## 🧪 Test Rapide

### 1. Créer un compte SuperAdmin
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Admin",
    "email": "admin@test.com",
    "password": "secret123",
    "role": "SuperAdmin",
    "shop_name": "Ma Boutique",
    "whatsapp_number": "212612345678"
  }'
```

### 2. Se connecter
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "secret123"
  }'
```

### 3. Tester la route publique
```bash
curl http://localhost:8080/public/shops
```

---

## 📊 Architecture du Projet
```
electronic-shop-api/
├── config/
│   └── config.go           # Configuration JWT & serveur
├── database/
│   └── database.go         # Connexion SQLite + migrations
├── handlers/
│   ├── auth.go             # Register, Login, GetMe
│   ├── products.go         # CRUD Produits + Routes publiques
│   ├── transactions.go     # CRUD Transactions
│   ├── dashboard.go        # Dashboard & Rapports
│   └── shop.go             # Gestion Shop & Utilisateurs
├── middleware/
│   └── auth.go             # JWT Middleware + Rôles
├── models/
│   └── models.go           # Shop, User, Product, Transaction
├── routes/
│   └── routes.go           # Configuration des routes
├── frontend/
│   ├── index.html          # Page principale
│   ├── style.css           # Styles
│   └── app.js              # Logique JavaScript
├── main.go                 # Point d'entrée
├── Dockerfile              # Image Docker
├── docker-compose.yml      # Orchestration
├── go.mod                  # Dépendances Go
└── README.md               # Ce fichier
```

---

## 🌐 Endpoints API

### 🔓 Routes Publiques (Sans authentification)

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| POST | `/register` | Créer un compte |
| POST | `/login` | Se connecter |
| GET | `/public/shops` | Liste des shops actifs |
| GET | `/public/:shopID/products` | Produits d'un shop |
| GET | `/public/:shopID/products/:id` | Détail produit + lien WhatsApp |

### 🔐 Routes Protégées (JWT requis)

| Méthode | Endpoint | Rôle | Description |
|---------|----------|------|-------------|
| GET | `/me` | Tous | Profil utilisateur |
| GET | `/products` | Admin+ | Liste des produits |
| POST | `/products` | Admin+ | Créer un produit |
| PUT | `/products/:id` | Admin+ | Modifier un produit |
| DELETE | `/products/:id` | Admin+ | Supprimer un produit |
| GET | `/transactions` | Admin+ | Liste des transactions |
| POST | `/transactions` | Admin+ | Créer une transaction |
| DELETE | `/transactions/:id` | Admin+ | Supprimer une transaction |
| GET | `/reports/dashboard` | SuperAdmin | Dashboard complet |
| GET | `/reports/low-stock` | SuperAdmin | Produits stock faible |
| GET | `/shop` | SuperAdmin | Info du shop |
| PUT | `/shop` | SuperAdmin | Modifier le shop |
| GET | `/users` | SuperAdmin | Liste des utilisateurs |
| POST | `/users` | SuperAdmin | Créer un utilisateur |
| PUT | `/users/:id` | SuperAdmin | Modifier un utilisateur |
| DELETE | `/users/:id` | SuperAdmin | Supprimer un utilisateur |

---

## 🔐 Rôles & Permissions

| Permission | SuperAdmin | Admin | Guest |
|------------|:----------:|:-----:|:-----:|
| Voir produits publics | ✅ | ✅ | ✅ |
| CRUD Produits | ✅ | ✅ | ❌ |
| Voir PurchasePrice | ✅ | ❌ | ❌ |
| CRUD Transactions | ✅ | ✅ | ❌ |
| Voir Dashboard/Profits | ✅ | ❌ | ❌ |
| Gérer Utilisateurs | ✅ | ❌ | ❌ |
| Modifier WhatsApp | ✅ | ❌ | ❌ |

---

## 🔒 Sécurité Multi-tenant

- ✅ Chaque requête est filtrée par `ShopID` extrait du JWT
- ✅ Un utilisateur ne peut JAMAIS accéder aux données d'un autre shop
- ✅ `PurchasePrice` n'est JAMAIS exposé au public ou aux Admin
- ✅ Mots de passe hashés avec bcrypt
- ✅ Tokens JWT avec expiration (24h)

---

## 📱 Intégration WhatsApp

Les routes publiques génèrent automatiquement un lien WhatsApp :
```
https://wa.me/212612345678?text=Bonjour%20je%20veux%20plus%20d'information%20sur%20iPhone%2014
```

---

## 🐛 Dépannage

### Le backend ne démarre pas ?
```bash
# Vérifier que le port 8080 est libre
lsof -i :8080

# Relancer avec les logs
go run main.go 2>&1
```

### Erreur CORS ?

Le middleware CORS est déjà configuré dans `main.go`. Vérifiez que vous utilisez bien `http://localhost:8080` comme URL de l'API.

### Erreur JWT ?

Vérifiez que vous envoyez bien le header :
```
Authorization: Bearer <votre-token>
```

---

## 📝 Licence

Projet éducatif - Bootcamp Go 2024

---

## 🙏 Remerciements

Merci à **Mr. Mounir Aziz** pour son encadrement tout au long de ce bootcamp.