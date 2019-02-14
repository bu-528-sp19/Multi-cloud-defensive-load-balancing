document.querySelector(".container").style.display = 'block';



document.querySelector('#submitTask').addEventListener('click', decideMethod);


function decideMethod() {
    if (document.querySelector("#args").value == "") {
        getRequest();
    }
    else {
        postRequest();
    }
}


function getRequest() {
        let j = 0;
        while(document.querySelector(`#user${j}`).value != "") {
            document.querySelector(`#user${j}`).value = ""
            j++;
        }

        let param = document.querySelector("#method").value;
        let id = document.querySelector("#args").value;

        //if id == nil { id = ""; }

        let request = new XMLHttpRequest();
        request.open('GET', `${param}`, true);
        request.onload = function () {
            let data = JSON.parse(this.response);
            if (data.length) {
                let i = 0;
                while(i<data.length) {
                    document.querySelector(`#user${i}`).value = JSON.stringify(data[i]);
                    i++;
                }
            }
            else {document.querySelector(`#user${0}`).value = JSON.stringify(data);}
            
            document.querySelector("#method").value = "";
            document.querySelector("#args").value = "";
        }
    request.send()
}


function postRequest() {
    let j = 0;
    while(document.querySelector(`#user${j}`).value != "") {
        document.querySelector(`#user${j}`).value = "";
        j++;
    }
    let param = document.querySelector("#method").value;
    let jsonobj = document.querySelector("#args").value;

    //if id == nil { id = ""; }

        let request = new XMLHttpRequest();
        request.open('POST', `${param}`, true);
        
        /*request.onload = function () {
            document.querySelector("#method").value = "";
            document.querySelector("#args").value = "";
        }*/
        console.log(jsonobj);

        request.send(jsonobj);
        
        document.querySelector("#method").value = "";
        document.querySelector("#args").value = "";

}


