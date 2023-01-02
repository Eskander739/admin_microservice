package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

const adminId = "5f5130c4-74ef-3f13-af4c-0b5137a36fe8"
const serverDomen = "http://localhost:8050"

func returnS(name string, w http.ResponseWriter) {
	//указываем путь к нужному файлу
	path := filepath.Join("static", name)
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func returnSEdit(name string, w http.ResponseWriter, data AddData) {
	path := filepath.Join("static", name)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

type IdImage struct {
	Id string `json:"id"`
}

type AddData struct {
	TitleForAdd   string
	ContentForAdd string
	IdForAdd      string
}

func main() {
	InitImagesDbAdmin()
	InitPostsDbAdmin()

	var editData = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				postId := mux.Vars(r)["post-id"]

				if postId != "" {
					var postParams = map[string]string{"Title": "TEXT", "Content": "TEXT"}

					var tables = map[string]map[string]string{"posts": postParams}
					var posts = Db{DbName: "eska", TableName: "posts", FetchInfo: "posts", Tables: tables}

					var data1, _ = posts.fetchInfo()
					var dataToSend PostData

					for _, info := range data1 {
						if info.(PostData).Id == postId {
							dataToSend.Title = info.(PostData).Title
							dataToSend.Content = info.(PostData).Content

						}
					}

					var data = AddData{TitleForAdd: dataToSend.Title, ContentForAdd: dataToSend.Content, IdForAdd: postId}
					returnSEdit("edit_post.html", w, data)
				}

			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var postList = func(w http.ResponseWriter, r *http.Request) {
		var posts = Db{DbName: "eska", TableName: "posts", FetchInfo: "posts"}
		var data1, _ = posts.fetchInfo()
		var data = map[string]any{"posts": data1}
		var dataPosts, jsonError = json.MarshalIndent(data, "", "   ")
		if jsonError != nil {
			panic(jsonError)
		}
		_, err := w.Write(dataPosts)
		if err != nil {
			panic(err)
		}

	}

	var saveChanges = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {

				decoder, errReadAll := ioutil.ReadAll(r.Body)
				if errReadAll != nil {
					panic(errReadAll)
				}

				var data PostData

				err := json.Unmarshal(decoder, &data)
				if err != nil {
					panic(err)
				}
				var postParams = map[string]string{"Title": "TEXT", "Content": "TEXT"}

				var tables = map[string]map[string]string{"posts": postParams}
				var post = PostData{Id: data.Id, Title: data.Title, Content: data.Content}

				var posts = Db{DbName: "eska", TableName: "posts", PostD: post, FetchInfo: "posts", Tables: tables}
				err = posts.ChangePost()
				if err != nil {
					panic(err)
				}
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var deletePost = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				postId := mux.Vars(r)["post-id"]

				var post = PostData{Id: postId}

				var posts = Db{DbName: "eska", TableName: "posts", PostD: post, FetchInfo: "posts"}
				_, err := posts.removeInfo()
				if err != nil {
					panic(err)
				}

				var imageS = ImageServ{Id: postId}
				var images = Db{DbName: "eska", TableName: "post_images", ImageS: imageS}
				_, err = images.removeInfoImage()
				if err != nil {
					panic(err)
				}
				http.Redirect(w, r, serverDomen+"/admin/posts", http.StatusSeeOther)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var addPostData = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				_, data, imageErr := r.FormFile("image")
				if imageErr != nil {
					panic(imageErr)
				}
				fileContent, err := data.Open()
				if err != nil {
					panic(err)
				}
				image, err := ioutil.ReadAll(fileContent)
				if err != nil {
					panic(err)
				}

				var uuidSQL = uuid4SQL()

				date := fmt.Sprintln(time.Now().Date())

				var post = PostData{Id: uuidSQL, Title: r.FormValue("title"), Content: r.FormValue("content"), Date: date}

				var imageData = ImageServ{Id: uuidSQL, Image: image}
				var posts = Db{DbName: "eska", TableName: "posts", PostD: post, FetchInfo: "posts", ImageS: imageData}
				err = posts.AddPost()
				if err != nil {
					panic(err)
				}

				err = posts.AddImage()
				if err != nil {
					panic(err)
				}

				returnS("add_post.html", w)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var add_post = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				returnS("add_post.html", w)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}
	}

	var posts = func(w http.ResponseWriter, r *http.Request) {

		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				returnS("all_posts.html", w)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var postPage = func(w http.ResponseWriter, r *http.Request) {
		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				returnS("posts.html", w)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var blog = func(w http.ResponseWriter, r *http.Request) {
		var cookie, err = r.Cookie("EskaUser")
		if err != nil {
			http.Redirect(w, r, serverDomen, http.StatusSeeOther)
		} else {

			if cookie.Value == adminId {
				returnS("blog_management.html", w)
			} else {
				w.WriteHeader(403)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)

			}

		}

	}

	var signInAdmin = func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		password2 := mux.Vars(r)["password"]

		if name == "EskanderAdminbdaea059c142ad7c463" {

			if password2 == "11736ad59330cbdaea059c142ad7c463" {
				cookie1 := http.Cookie{Name: "EskaUser", Value: adminId, Expires: time.Now().Add(time.Hour), HttpOnly: false, MaxAge: 50000, Path: "/"}
				http.SetCookie(w, &cookie1)
				http.Redirect(w, r, serverDomen, http.StatusSeeOther)
			}

		}

	}

	var getImageByIdAdmin = func(w http.ResponseWriter, r *http.Request) {
		decoder, errReadAll := ioutil.ReadAll(r.Body)
		if errReadAll != nil {
			panic(errReadAll)
		}
		var data IdImage

		err := json.Unmarshal(decoder, &data)
		if err != nil {
			panic(err)
		}
		var imageId = ImageServ{Id: data.Id}

		var images = Db{DbName: "eska", TableName: "post_images", FetchInfo: "post_images", ImageS: imageId}
		var data1, err2 = images.getImageById()
		if err2 != nil {
			panic(err2)
		}
		_, err = w.Write(data1)
		if err != nil {
			panic(err)
		}

	}

	var getPostByIdAdmin = func(w http.ResponseWriter, r *http.Request) {
		decoder, errReadAll := ioutil.ReadAll(r.Body)
		if errReadAll != nil {
			panic(errReadAll)
		}
		var data PostData

		err := json.Unmarshal(decoder, &data)
		if err != nil {
			panic(err)
		}
		var imageId = ImageServ{Id: data.Id}

		var images = Db{DbName: "eska", TableName: "posts", FetchInfo: "posts", ImageS: imageId, PostD: data}
		var data1, err2 = images.getPostById()
		if err2 != nil {
			panic(err2)
		}
		_, err = w.Write(data1)
		if err != nil {
			panic(err)
		}

	}

	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css/"))))          // ПОДКЛЮЧАЕМ CSS ФАЙЛЫ
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img/"))))          // ПОДКЛЮЧАЕМ ИЗОБРАЖЕНИЯ ФАЙЛЫ
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js/"))))             // ПОДКЛЮЧАЕМ JS ФАЙЛЫ
	router.PathPrefix("/vendor/").Handler(http.StripPrefix("/vendor/", http.FileServer(http.Dir("./static/vendor/")))) // ПОДКЛЮЧАЕМ VENDOR ФАЙЛЫ
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("./static/fonts/"))))    // ПОДКЛЮЧАЕМ FONTS ФАЙЛЫ

	// ***********************************ADMIN SECTION*********************************************
	router.HandleFunc("/admin", blog)                               // защита есть
	router.HandleFunc("/admin/post-page", postPage)                 // защита есть
	router.HandleFunc("/admin/new-post-page", add_post)             // защита есть
	router.HandleFunc("/admin/new-post", addPostData)               // защита есть
	router.HandleFunc("/admin/posts", posts)                        // защита есть
	router.HandleFunc("/post-list", postList)                       // защита есть
	router.HandleFunc("/admin/posts/{post-id}/changed", editData)   // защита есть
	router.HandleFunc("/admin/modified-post", saveChanges)          // защита есть
	router.HandleFunc("/admin/posts/{post-id}/deleted", deletePost) // защита есть

	router.HandleFunc("/admin/sign-in/{name}/{password}", signInAdmin) // защита есть
	// ***********************************ADMIN SECTION*********************************************

	router.HandleFunc("/get-image-by-id", getImageByIdAdmin)
	router.HandleFunc("/get-post-by-id", getPostByIdAdmin)

	listenError := http.ListenAndServe(":8050", router)

	if listenError != nil {
		panic(listenError)
	}
}

//http://localhost:8050/sign-in/EskanderAdminbdaea059c142ad7c463/11736ad59330cbdaea059c142ad7c463
