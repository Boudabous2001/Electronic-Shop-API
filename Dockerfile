# Dockerfile
FROM golang:1.22-alpine

# Installer les dépendances système
RUN apk add --no-cache gcc musl-dev

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
RUN go build -o main .

# Exposer le port
EXPOSE 8080

# Lancer l'application
CMD ["./main"]