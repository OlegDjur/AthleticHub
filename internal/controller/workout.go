package controller

// func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
// 	var workout model.Workout
// 	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	workout.CreatedAt = time.Now()
// 	workout.UpdatedAt = time.Now()

// 	if err := h.workoutService.CreateWorkout(r.Context(), &workout); err != nil {
// 		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(workout)
// }

// func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
// 		return
// 	}

// 	workout, err := h.workoutService.GetWorkout(r.Context(), id)
// 	if err != nil {
// 		http.Error(w, "Workout not found", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(workout)
// }

// func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
// 		return
// 	}

// 	var workout model.Workout
// 	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	workout.ID = id
// 	workout.UpdatedAt = time.Now()

// 	if err := h.workoutService.UpdateWorkout(r.Context(), &workout); err != nil {
// 		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(workout)
// }

// func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.workoutService.DeleteWorkout(r.Context(), id); err != nil {
// 		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
