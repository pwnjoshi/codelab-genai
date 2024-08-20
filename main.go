package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/vertexai/genai"
)

func main() {
	ctx := context.Background()
	var projectId string
	var err error
	projectId = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		projectId, err = metadata.ProjectIDWithContext(ctx)
		if err != nil {
			return
		}
	}
	var client *genai.Client
	client, err = genai.NewClient(ctx, projectId, "us-central1")
	if err != nil {
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash-001")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		animal := r.URL.Query().Get("animal")
		if animal == "" {
			animal = "dog"
		}

		resp, err := model.GenerateContent(
			ctx,
			genai.Text(
				fmt.Sprintf("Give me 10 fun facts about %s. Return the results as HTML without markdown backticks.", animal)),
		)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			htmlContent := resp.Candidates[0].Content.Parts[0]
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			
			// Adding basic HTML structure and inline CSS for styling
			fmt.Fprintf(w, `
				<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Fun Facts about %s</title>
					<style>
						body {
							font-family: Arial, sans-serif;
							background-color: #f0f8ff;
							margin: 0;
							padding: 0;
							display: flex;
							justify-content: center;
							align-items: center;
							height: 100vh;
							color: #333;
						}
						.container {
							background-color: #fff;
							padding: 20px;
							border-radius: 8px;
							box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
							max-width: 600px;
							width: 100%;
							text-align: left;
						}
						h1 {
							font-size: 24px;
							color: #333;
							margin-bottom: 20px;
						}
						p {
							font-size: 18px;
							line-height: 1.6;
						}
					</style>
				</head>
				<body>
					<div class="container">
						<h1>10 Fun Facts about %s</h1>
						%s
					</div>
				</body>
				</html>
			`, animal, animal, htmlContent)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
