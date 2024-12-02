// src/store/reducers/index.js
import { combineReducers } from 'redux';
import authReducer from './authReducer';
// Добавьте другие редьюсеры по мере необходимости

const rootReducer = combineReducers({
    auth: authReducer,
    // other reducers
});

export default rootReducer;
