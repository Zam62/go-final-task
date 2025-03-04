document.addEventListener('DOMContentLoaded', () => {
    const expressionInput = document.getElementById('expression');
    const fetchDataBtn = document.getElementById('fetchdata');
    const resultDiv = document.getElementById('result');
    const historyList = document.getElementById('history-list');


    fetchDataBtn.addEventListener('click', () => {
        calculateExpression()
    });



    async function calculateExpression() {

        alert('go');
        result.innerText = 'Loading....'

        const expression = expressionInput.value.trim();
        if (!expression) {
            result.innerText = 'Bad expression....'
            return;
        }
        
        try {
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

        } catch (error) {
            showError(error.message);
        }
    }    
    });




const host = 'http://api/v1/calculate';

const fetchData = function() {
    alert('go');
result.innerText = 'Loading....'
fetch(host)
    .then(res => res.json())
    .then(data => {
    result.innerText = JSON.stringify(data, null, 2)
    })
    .catch(error => console.log(error))
}
