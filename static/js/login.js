const selector = (e) => document.querySelector(e);

const sendData = async (url, data) => {
    try {
        let useData = await fetch(url, {
            method: "POST",
            redirect: "follow",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
        let parseRes = await useData.json();
        return parseRes;
    } catch (error) {
        return error;
    }
};

let sendForm = selector("#formSubmit");
let serverResponse = selector("#response");

sendForm.addEventListener("click", (e) => {
    e.preventDefault();

    sendData("/login", {
        email: selector("#email").value,
        password: selector("#password").value,
    })
        .then((res) => {
            console.log(res);
            serverResponse.textContent = "";
            serverResponse.textContent = `Your email is ${res.email}`;
        })
        .catch((err) => (serverResponse.textContent = err.message));
});
