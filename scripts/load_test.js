import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    stages: [
        { duration: '30s', target: 300 }, // Разогрев до 300 пользователей
        { duration: '2m', target: 300 },  // Стабильная нагрузка 2 минуты
        { duration: '30s', target: 0 },   // Снижение нагрузки
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% запросов должны быть быстрее 500мс
        http_req_failed: ['rate<0.01'],   // Ошибок должно быть меньше 1%
    },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
    const email = `user_${randomString(10)}@example.com`;
    const password = 'StrongPassword123!';

    // 1. Регистрация
    const registerRes = http.post(`${BASE_URL}/auth/register`, JSON.stringify({
        email: email,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });

    check(registerRes, {
        'status is 201 (registered)': (r) => r.status === 201,
    });

    sleep(0.1);

    // 2. Логин
    const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
        email: email,
        password: password,
    }), { headers: { 'Content-Type': 'application/json' } });

    const loginOk = check(loginRes, {
        'status is 200 (logged in)': (r) => r.status === 200,
        'has tokens': (r) => r.json().access_token !== undefined && r.json().refresh_token !== undefined,
    });

    if (!loginOk) return;

    const accessToken = loginRes.json().access_token;
    const refreshToken = loginRes.json().refresh_token;

    sleep(0.1);

    // 3. Получение своего профиля (Protected Route)
    const meRes = http.get(`${BASE_URL}/accounts/me`, {
        headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
        },
    });

    check(meRes, {
        'status is 200 (me profile)': (r) => r.status === 200,
        'email matches': (r) => r.json().email === email,
    });

    sleep(0.1);

    // 4. Рефреш токена
    const refreshRes = http.post(`${BASE_URL}/auth/refresh`, JSON.stringify({
        refresh_token: refreshToken,
    }), { headers: { 'Content-Type': 'application/json' } });

    check(refreshRes, {
        'status is 200 (token refreshed)': (r) => r.status === 200,
        'new tokens generated': (r) => r.json().access_token !== undefined,
    });

    const newRefreshToken = refreshRes.json().refresh_token;

    sleep(0.1);

    // 5. Логаут
    const logoutRes = http.post(`${BASE_URL}/auth/logout`, JSON.stringify({
        refresh_token: newRefreshToken,
    }), { headers: { 'Content-Type': 'application/json' } });

    check(logoutRes, {
        'status is 200 (logged out)': (r) => r.status === 200,
    });

    sleep(1);
}
