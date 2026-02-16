# 📄 Workflow IA - Electronic Shop API

## Bootcamp Go 2024

**Équipe** : Boudabous Elyes, Yasmine Aoudjit, Mohamed Amine Dhaoui, Mohamed Amine Ourraki

**Encadrant** : Mr. Mounir Aziz

---

## 1️⃣ Méthode de Travail

### Architecture choisie

Nous avons opté pour une architecture **MVC (Model-View-Controller)** adaptée à Go :
```
├── models/      → Modèles de données (GORM)
├── handlers/    → Contrôleurs (logique métier)
├── routes/      → Définition des endpoints
├── middleware/  → JWT & Vérification des rôles
```

### Planification

1. **Phase 1** : Définition des modèles de données
2. **Phase 2** : Authentification JWT
3. **Phase 3** : CRUD Produits & Transactions
4. **Phase 4** : Dashboard & Rapports
5. **Phase 5** : Routes publiques & WhatsApp
6. **Phase 6** : Frontend
7. **Phase 7** : Tests & Documentation

---

## 2️⃣ Outils IA Utilisés

| Outil | Usage |
|-------|-------|
| **Claude (Anthropic)** | Génération de code, architecture, débogage |
| **GitHub Copilot** | Autocomplétion dans VS Code |

---

## 3️⃣ Prompts Importants

### Prompt 1 : Structure du projet
```
"Crée la structure d'un projet Go pour une API REST de gestion 
de boutiques d'électronique avec :
- Authentification JWT
- Multi-tenant (isolation par shop)
- Rôles SuperAdmin et Admin
- GORM avec SQLite"
```

### Prompt 2 : Middleware d'authentification
```
"Crée un middleware JWT en Go avec Gin qui :
- Extrait le token du header Authorization
- Valide le token
- Stocke userID, shopID et role dans le contexte
- Permet de vérifier les rôles requis"
```

### Prompt 3 : Routes publiques
```
"Crée une route publique GET /public/:shopID/products qui :
- Ne nécessite pas d'authentification
- Retourne les produits sans le PurchasePrice
- Génère un lien WhatsApp dynamique pour chaque produit"
```

### Prompt 4 : Correction CGO/SQLite
```
"J'ai l'erreur 'Binary was compiled with CGO_ENABLED=0, go-sqlite3 
requires cgo' sur Windows. Comment résoudre sans installer CGO ?"
```

**Solution obtenue** : Utiliser `github.com/glebarez/sqlite` (driver pure Go)

---

## 4️⃣ Intégration Front / Back

### Connexion API
```javascript
const API_URL = 'http://localhost:8080';

// Appel avec token JWT
const res = await fetch(`${API_URL}/products`, {
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    }
});
```

### Gestion du Token

- Stockage dans `localStorage`
- Envoi automatique dans le header `Authorization`
- Vérification à chaque chargement de page

### Gestion des Erreurs
```javascript
if (!res.ok) {
    const data = await res.json();
    alert(data.error); // Message d'erreur du backend
}
```

---

## 5️⃣ Difficultés Rencontrées

### 1. Problème CGO/SQLite sur Windows

**Erreur** : `go-sqlite3 requires cgo to work`

**Cause** : Le driver `gorm.io/driver/sqlite` nécessite CGO

**Solution** : Utiliser `github.com/glebarez/sqlite` (pure Go)

### 2. Problème CORS

**Erreur** : `Access-Control-Allow-Origin` bloqué

**Solution** : Ajout d'un middleware CORS dans `main.go`
```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        // ...
    }
}
```

### 3. PurchasePrice exposé aux Admin

**Problème** : L'IA générait du code qui exposait `purchase_price` à tous

**Solution** : Vérification manuelle du rôle avant de retourner les données
```go
if role == models.RoleAdmin {
    // Retourner sans purchase_price
}
```

### 4. Stock négatif possible

**Problème** : Pas de vérification du stock avant vente

**Solution** : Ajout de validation + transaction atomique
```go
if product.Stock < input.Quantity {
    return error("Stock insuffisant")
}
tx := db.Begin()
// ...
tx.Commit()
```

---

## 6️⃣ Points Bloquants

### Ce qui n'a pas fonctionné

| Problème | Cause | Solution |
|----------|-------|----------|
| SQLite sur Windows | CGO requis | Driver pure Go |
| Token non envoyé | Oubli du préfixe "Bearer " | Correction du format |
| Multi-tenant cassé | Filtrage par ShopID manquant | Ajout systématique |

### Ce que l'IA a mal généré

1. **Oubli du filtre ShopID** dans certaines requêtes
2. **Exposition du PurchasePrice** dans les routes publiques
3. **Pas de validation** avant décrémentation du stock
4. **Messages d'erreur génériques** au lieu de messages précis

---

## 7️⃣ Analyse Critique de l'IA

### ✅ Où l'IA nous a fait gagner du temps

| Tâche | Sans IA | Avec IA | Gain |
|-------|---------|---------|------|
| Structure projet | 3h | 30min | 83% |
| Modèles GORM | 2h | 15min | 87% |
| Middleware JWT | 3h | 30min | 83% |
| Handlers CRUD | 5h | 1h | 80% |
| Frontend HTML/CSS | 4h | 1h | 75% |
| Documentation | 2h | 30min | 75% |
| **Total** | **19h** | **~4h** | **79%** |

### ❌ Où l'IA nous a fait perdre du temps

1. **Débogage CGO** : 1h30 à chercher pourquoi ça ne compilait pas
2. **Sécurité multi-tenant** : 1h à vérifier et corriger toutes les requêtes
3. **Tests manuels** : 2h pour vérifier que tout fonctionne correctement

### 🔧 Corrections manuelles apportées

1. Remplacement du driver SQLite
2. Ajout du filtre `shop_id` dans TOUTES les requêtes
3. Masquage du `purchase_price` pour les Admin
4. Validation du stock avant vente
5. Messages d'erreur en français
6. Protection contre l'auto-suppression

---

## 8️⃣ Conclusion

L'IA est un **accélérateur puissant** mais elle nécessite :

- ✅ Une supervision humaine constante
- ✅ Des tests manuels de sécurité
- ✅ Une compréhension approfondie du code généré
- ✅ Des corrections et améliorations manuelles

Le développeur reste **responsable** de la qualité et de la sécurité du code final.

---

## 📊 Récapitulatif

| Critère | Évaluation |
|---------|------------|
| Temps gagné | ~79% |
| Corrections nécessaires | ~15 modifications |
| Bugs générés par l'IA | 4 bugs majeurs |
| Qualité du code initial | 7/10 |
| Qualité après corrections | 9/10 |

---

*Document rédigé par l'équipe dans le cadre du Bootcamp Go 2024*