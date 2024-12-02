import axios from 'axios';

// Установить токен в заголовки
const setAuthToken = (token) => {
    if (token) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        delete axios.defaults.headers.common['Authorization'];
    }
};

export const loadUser = () => async (dispatch) => {
    const token = localStorage.getItem('token');
    if (token) {
        setAuthToken(token); // Установить токен
    } else {
        return dispatch({ type: 'AUTH_ERROR' });
    }

    try {
        const res = await axios.get('/api/users/me');
        dispatch({
            type: 'USER_LOADED',
            payload: res.data, // Данные пользователя
        });
        console.log('User loaded:', res.data);
    } catch (err) {
        console.error(err);
        dispatch({
            type: 'AUTH_ERROR',
        });
        setAuthToken(null); // Удалить токен из заголовков
        localStorage.removeItem('token'); // Очистить токен
    }
};

export const register = (userData) => async (dispatch) => {
    try {
        const res = await axios.post('/api/register', userData);
        dispatch({
            type: 'REGISTER_SUCCESS',
            payload: res.data, // Токен и данные пользователя
        });

        const { token } = res.data;
        localStorage.setItem('token', token);
        setAuthToken(token);

        // Загружаем пользователя после регистрации
        dispatch(loadUser());
    } catch (err) {
        console.error('Registration error:', err);
        dispatch({
            type: 'REGISTER_FAIL',
            payload: err.response?.data?.message || 'Registration failed',
        });
        throw err; // Передать ошибку в компонент
    }
};

export const login = (credentials) => async (dispatch) => {
    try {
        const res = await axios.post('/api/login', credentials);
        dispatch({
            type: 'LOGIN_SUCCESS',
            payload: res.data, // Токен и данные пользователя
        });

        const { token } = res.data;
        localStorage.setItem('token', token);
        setAuthToken(token);

        // Загружаем пользователя после логина
        dispatch(loadUser());
    } catch (err) {
        console.error('Login error:', err);
        dispatch({
            type: 'LOGIN_FAIL',
            payload: err.response?.data?.message || 'Login failed',
        });
        throw err; // Передать ошибку в компонент
    }
};

export const logout = () => (dispatch) => {
    localStorage.removeItem('token');
    setAuthToken(null);
    dispatch({ type: 'LOGOUT' });
};
