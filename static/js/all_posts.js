var serverDomen = "http://localhost:8050"

function getData(callback) {
    var xhr = new XMLHttpRequest();
    // we defined the xhr
    var yourUrl = serverDomen + "/post-list";

    xhr.onreadystatechange = function () {
        if (this.readyState == 4) {
            callback.apply(xhr);

        }

        if (this.status == 200) {
            var data = JSON.parse(this.responseText);



            // we get the returned data
        }


        // end of state change: it can be after some time (async)

    };


    xhr.open('POST', yourUrl, true);
    xhr.send();

    // return data1;


}

function decoderImg(buffer) {
    alert(new Uint8Array(buffer));
    var mime;
    var a = new Uint8Array(b(buffer));
    var nb = a.length;
    if (nb < 4)
        return null;
    alert("lengrh");
    var b0 = a[0];
    var b1 = a[1];
    var b2 = a[2];
    var b3 = a[3];
    if (b0 == 0x89 && b1 == 0x50 && b2 == 0x4E && b3 == 0x47)
        mime = 'image/png';
    else if (b0 == 0xff && b1 == 0xd8)
        mime = 'image/jpeg';
    else if (b0 == 0x47 && b1 == 0x49 && b2 == 0x46)
        mime = 'image/gif';
    else
        return null;
    alert("eeee");
    var binary = "";
    for (var i = 0; i < nb; i++)
        binary += String.fromCharCode(a[i]);
    var base64 = window.btoa(binary);

    return 'data:' + mime + ';base64,' + base64;
}

async function getImage(imageId) {
    // alert(1);
    let xhr = new XMLHttpRequest();
    let url = serverDomen + "/get-image-by-id";

    xhr.onreadystatechange = function () {
        if (this.readyState == 4) {

        }

        if (this.status == 200) {
            var immgg = document.getElementById(imageId);
            var blob = new Blob([this.response], { type: "image/jpg" });
            immgg.src = URL.createObjectURL(blob);

        }


        // end of state change: it can be after some time (async)

    };
    xhr.open("POST", url, true);
    xhr.responseType = "arraybuffer"; //ОБЯЗАТЕЛЬНЫЙ ТЕГ ЧТОБЫ ИЗОБРАЖЕНИЕ ЧИТАЛОСЬ!!!
    xhr.setRequestHeader("Content-Type", "application/json");
    var data = JSON.stringify({ "id": imageId });
    xhr.send(data);
    // alert(xhr.response);


}

function DataPosts() {

    var resp = JSON.parse(this.responseText);
    for (index = 0; index < resp.posts.length; ++index) {
        post = resp.posts[index];
        var last_id = main_div.children.length;
        var last_elem = window.getComputedStyle(main_div.children[last_id - 1]).left;
        var h1Data = document.createElement('h1');
        h1Data.className = "zagolovok";
        var title = document.createElement('label');

        var postImage = document.createElement('img');

        var buttonLinkB = document.createElement('a');
        buttonLinkB.className = "nav-link";
        buttonLinkB.innerText = "Удалить";
        buttonLinkB.href = serverDomen + `/admin/posts/${post.Id}/deleted`;

        var buttonLinkA = document.createElement('a');
        buttonLinkA.className = "nav-link";
        buttonLinkA.innerText = "Редактировать";
        buttonLinkA.href = serverDomen + `/admin/posts/${post.Id}/changed`;//${post.title}/${post.content}
        var buttonRemove = document.createElement('button');
        var buttonEdit = document.createElement('button');
        var idd = document.createElement('label');
        idd.className = "idForAdd";
        idd.innerText = post.Id;
        idd.style = "display: none";
        buttonRemove.className = "remove-post";
        buttonRemove.appendChild(buttonLinkB);
        buttonEdit.className = "edit-post";
        buttonEdit.appendChild(buttonLinkA);
        // buttonEdit.innerText = "Редактировать";
        postImage.id = post.Id;
        postImage.style = "display: block;width: auto;height: auto;max-width: 100%;";



        h1Data.innerText = post.title;
        title.className = 'title';
        title.appendChild(h1Data);
        var content = document.createElement('label');
        content.className = "content";
        content.innerText = post.content;
        content.style = "overflow: hidden;display: -webkit-box;-webkit-line-clamp: 3;-webkit-box-orient: vertical"
        var menuItem = document.createElement('div');
        menuItem.className = 'menu--item';
        var iDiv = document.createElement('figure');
        
        // iDiv.style.left = `300px`;
        // preee.appendChild(content);
        menuItem.style = `

  
        box-shadow:
         0 1px 4px rgba(0, 0, 0, .3),
         -23px 0 20px -23px rgba(0, 0, 0, .8),
         23px 0 20px -23px rgba(0, 0, 0, .8),
         0 0 40px rgba(0, 0, 0, .1) inset;
      }`;
        menuItem.style.left = `${Number(last_elem.replace("px", "")) + 460}px`;

        iDiv.appendChild(title);
        iDiv.appendChild(buttonRemove);
        iDiv.appendChild(buttonEdit);
        iDiv.appendChild(content);
        iDiv.appendChild(postImage);
        iDiv.appendChild(idd);
        // iDiv.id = 'block';

        
        menuItem.appendChild(iDiv);
        
        document.getElementsByClassName('menu--wrapper')[0].appendChild(menuItem);
        getImage(post.Id);
    }

}
function addPost() {
    getData(DataPosts);
}
