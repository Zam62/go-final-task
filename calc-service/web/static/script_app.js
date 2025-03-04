document.addEventListener('DOMContentLoaded', () => {
    const expressionInput = document.getElementById('expression');
    const fetchDataBtn = document.getElementById('fetchdata');
    const resultDiv = document.getElementById('result');
    const historyList = document.getElementById('history-list');

    // loadHistory();

    fetchDataBtn.addEventListener('click', () => {
        // getData();
        calculateExpression()
    });



    async function calculateExpression() {

        alert('go');
        result.innerText = 'Loading....'

        const expression = expressionInput.value.trim();
        if (!expression) {
            result.innerText = 'Bad expression....'
            // showError('Пожалуйста, введите выражение');
            return;
        }
        
        try {
            // resultDiv.innerHTML = '<div class="processing">Вычисление...</div>';
            result.innerText = 'Calculating....'

            const response = await fetch('/api/v1/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ expression })
            });

            result.innerText = response.statusText
            
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Ошибка: ${response.status} ${response.statusText}`);
            }
            
            const data = await response.json();
            result.innerText = JSON.stringify(data, null, 2)
            const expressionId = data.id;

            // checkExpressionStatus(expressionId);
        } catch (error) {
            showError(error.message);
        }
    }    


    // expressionInput.addEventListener('keypress', (e) => {
    //     if (e.key === 'Enter') {
    //         calculateExpression();
    //     }
    });



// const fetchDataBtn = document.getElementById('fetchdata');
// // const fetchDataBtn = document.querySelector('#fetchdata')
// const result = document.querySelector('#result')

// alert('hi');

// // gets data from API and sets the content of #result div
// const host = 'http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/car_brand';
const host = 'http://api/v1/calculate';

const getData = function() {
    alert('go');
result.innerText = 'Loading....'
// fetch('https://dummyjson.com/products')
// fetch('http://api/v1/calculate')
fetch(host)
    .then(res => res.json())
    .then(data => {
    result.innerText = JSON.stringify(data, null, 2)
    })
    .catch(error => console.log(error))
}

// // add event listener for #fetchdata button
// fetchDataBtn.addEventListener('click', alert('gggg'))
