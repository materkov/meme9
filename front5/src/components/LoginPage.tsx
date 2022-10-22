import React from "react";
import styles from "./LoginPage.module.css";
import {api} from "../store/types";
import {Link} from "./Link";
import {authorize, navigate} from "../utils/localize";

const vkURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&redirect_uri=" + location.origin + "/vk-callback";

export function LoginPage() {
    const [email, setEmail] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [error, setError] = React.useState("");

    const onLogin = () => {
        setError('');

        api("/emailLogin", {
            email: email,
            password: password,
        }).then((resp) => {
            authorize(resp.token);
            navigate("/");
        }).catch(err => {
            if (err === 'invalid credentials') {
                setError('Неверный логин или пароль');
            }
        })
    }

    return <>
        <h2>Авторизация</h2>

        <input type="text" placeholder="Емейл" value={email} onChange={e => setEmail(e.target.value)}/>
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)}/>

        <button className={styles.submitBtn} onClick={onLogin}>Войти</button>

        {error &&
            <>
                <br/>
                <div className={styles.error}>{error}</div>
            </>
        }

        <br/><br/><br/>
        <a href={vkURL}>Войти через ВКонтакте</a>
        <br/>
        <Link href={"/register"}>Регистрация</Link>
    </>;
}
