const app = document.getElementById('root');

const logo = document.createElement('img');
logo.src = 'ebaumsworld.png';
const container = document.createElement('div');
container.setAttribute('class', 'container');

app.appendChild(logo);
app.appendChild(container);

function processData(response) {
    const data = JSON.parse(response).leaderboard.Users;

    if (document.contains(document.getElementById("wrapperElement"))) {
        document.getElementById("wrapperElement").remove();
    }

    const wrapper = document.createElement('div');
    wrapper.setAttribute("id", "wrapperElement");
    wrapper.setAttribute("class", "col-12");
    container.appendChild(wrapper);

    data.forEach(user => {
        const card = document.createElement('div');
        card.setAttribute('class', 'card');

        const h1 = document.createElement('h1');
        h1.textContent = `Player ${user.username}:${user.points}, is ranked #${user.rank}!`;

        const buttonBox = document.createElement('div');
        buttonBox.setAttribute('class', 'button-box');

        const buttonWin = document.createElement('button');
        buttonWin.setAttribute('class', 'btn btn-outline-dark col-3');
        buttonWin.textContent = "Add win"
        buttonWin.onclick = function(){incrementUser(user.username)}

        const buttonGetRival = document.createElement('button');
        buttonGetRival.setAttribute('class', 'btn btn-outline-dark col-3 offset-2');
        buttonGetRival.textContent = "Get rival"

        wrapper.appendChild(card);
        card.appendChild(h1);
        card.appendChild(buttonBox);
        buttonBox.appendChild(buttonWin);
        buttonBox.appendChild(buttonGetRival);
    });
}

function loadPage() {
    const request = new XMLHttpRequest();
    request.open('GET', 'http://localhost:8080/leaderboard', true);
    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            processData(this.response)
        } else {
            const errorMessage = document.createElement('marquee');
            errorMessage.textContent = `Gah, it's not working!`;
            app.appendChild(errorMessage);
        }
    }

    request.send();
}

function incrementUser(username) {
    const updateRequest = new XMLHttpRequest();
    updateRequest.open('POST', 'http://localhost:8080/points', true);
    updateRequest.setRequestHeader("Content-Type", "application/json")
    updateRequest.onload = function () {
        loadPage()
    }

    updateRequest.send(JSON.stringify({ "username": username}));

    // document.getElementById("demo").innerHTML = "YOU CLICKED ME!";
}

loadPage()
