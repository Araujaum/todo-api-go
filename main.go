package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Task representa uma tarefa com ID, título e status de conclusão.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// tasks é uma lista (slice) que armazena todas as tarefas.
var tasks []Task

// currentID é usado para gerar IDs únicos para cada tarefa.
var currentID = 1

func main() {
	// Configura os handlers para as rotas da API.
	http.HandleFunc("/tasks", handleTasks)       // Rota para listar e criar tarefas.
	http.HandleFunc("/tasks/", handleTask)       // Rota para atualizar e excluir tarefas.

	// Inicia o servidor na porta 8080.
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

// handleTasks gerencia as requisições para a rota "/tasks".
func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTasks(w, r) // Chama a função para listar tarefas.
	case "POST":
		createTask(w, r) // Chama a função para criar uma nova tarefa.
	default:
		// Retorna um erro se o método não for suportado.
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

// handleTask gerencia as requisições para a rota "/tasks/{id}".
func handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		updateTask(w, r) // Chama a função para atualizar uma tarefa.
	case "DELETE":
		deleteTask(w, r) // Chama a função para excluir uma tarefa.
	default:
		// Retorna um erro se o método não for suportado.
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

// getTasks retorna a lista de todas as tarefas em formato JSON.
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks) // Converte a lista de tarefas para JSON e envia como resposta.
}

// createTask cria uma nova tarefa com base nos dados enviados no corpo da requisição.
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	// Decodifica o JSON do corpo da requisição para a struct Task.
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atribui um ID único à tarefa e incrementa o currentID.
	task.ID = currentID
	currentID++
	// Adiciona a nova tarefa à lista de tarefas.
	tasks = append(tasks, task)

	// Retorna a tarefa criada em formato JSON com status 201 (Created).
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// updateTask atualiza uma tarefa existente com base no ID fornecido na URL.
func updateTask(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL (ex: "/tasks/1" -> ID = 1).
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da requisição para a struct Task.
	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Procura a tarefa pelo ID e atualiza seus dados.
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updatedTask.Title
			tasks[i].Completed = updatedTask.Completed
			// Retorna a tarefa atualizada em formato JSON.
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	// Retorna um erro se a tarefa não for encontrada.
	http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
}

// deleteTask exclui uma tarefa com base no ID fornecido na URL.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL (ex: "/tasks/1" -> ID = 1).
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Procura a tarefa pelo ID e a remove da lista.
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			// Retorna status 204 (No Content) para indicar sucesso sem conteúdo.
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Retorna um erro se a tarefa não for encontrada.
	http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
}

