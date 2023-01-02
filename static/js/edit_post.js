var serverDomen = "http://localhost:8050"
// function sendContentAndTitle(callback) {
//     var xhr = new XMLHttpRequest();
//     // we defined the xhr
//     var yourUrl = "http://localhost:8000/edit_data";

//     xhr.onreadystatechange = function () {
//         if (this.readyState == 4) {
//             callback.apply(xhr);

//         }

//         if (this.status == 200) {
//             var data = JSON.parse(this.responseText);

            

//             // we get the returned data
//         }


//         // end of state change: it can be after some time (async)

//     };


//     xhr.open('POST', yourUrl, true);
//     xhr.send();

//     // return data1;

    
// }


function DataPosts(){

    var resp  = JSON.parse(this.responseText);
    for (index = 0; index < resp.posts.length; ++index) {
        //<a class="nav-link" href="http://localhost:8000/all_posts"><i></i>Cancel</a>
        post = resp.posts[index];
        var last_id = main_div.children.length;
        var last_elem = window.getComputedStyle(main_div.children[last_id - 1]).left;
        var h1Data = document.createElement('h1');
        var title = document.createElement('label');
        var buttonLinkA = document.createElement('a');
        buttonLinkA.className = "nav-link";
        buttonLinkA.innerText = "Редактировать";
        buttonLinkA.href = serverDomen + `/edit_data/?id=${post.Id}`;//${post.title}/${post.content}
        var buttonRemove = document.createElement('button');
        var buttonEdit = document.createElement('button');

        buttonRemove.className = "remove-post";
        buttonRemove.innerText = "Удалить";
        buttonEdit.className = "edit-post";
        buttonEdit.appendChild(buttonLinkA);
        // buttonEdit.innerText = "Редактировать";

        h1Data.innerText = post.title;
        title.className = 'title';
        title.appendChild(h1Data);
        var content = document.createElement('label');
        content.innerText = post.content;
        var iDiv = document.createElement('div');
        iDiv.style.left = `${Number(last_elem.replace("px", "")) + 310}px`;
        // iDiv.style.left = `300px`;
        iDiv.className = 'post-type';
        iDiv.appendChild(title);
        iDiv.appendChild(content);
        iDiv.appendChild(buttonRemove);
        iDiv.appendChild(buttonEdit);
        // iDiv.id = 'block';

        document.getElementsByClassName('main_div')[0].appendChild(iDiv);

    }
    
}
function addPost() {
    getData();




    // alert(window.getComputedStyle(main_div.children[1]).left);
    // alert(main_div.children[last_id-1]);


}
