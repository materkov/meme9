import React, {useEffect} from "react";
import * as styles from "./Auth.module.css";
import {setAuth} from "../../store/globals";
import {navigationGo} from "../../store/navigation";
import {Link} from "../Link/Link";
import {cookieAuthToken, delCookie, setCookie} from "../../utils/cookie";
import {ApiAuth} from "../../api/client";

export function Auth() {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [isVK, setIsVK] = React.useState(false);

    const [error, setError] = React.useState('');
    const isReg = location.search === "?registration";

    // TODO think about location.search
    if (location.search === "?logout") {
        delCookie(cookieAuthToken);
        setAuth({token: "", userId: "", userName: ""});
        navigationGo("/");
        return null;
    }

    useEffect(() => {
        if (location.search.startsWith("?code")) {
            ApiAuth.Vk({
                code: new URLSearchParams(location.search).get('code') || "",
                redirectUrl: location.origin + location.pathname,
            }).then(resp => {
                setAuth(resp);

                setCookie(cookieAuthToken, resp.token);
                navigationGo("/");
            }).catch(() => {
                setIsVK(false);
                setError('Failed authorizing via VK');
            });

            setIsVK(true);
        }
    }, []);

    const vkAuthURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&v=5.131&redirect_uri=" + encodeURIComponent(location.origin + "/vk-callback");

    const onAuth = () => {
        if (!email || !password) {
            return;
        }
        setError('');

        if (isReg) {
            onRegister();
        } else {
            onLogin();
        }
    };

    const onLogin = () => {
        ApiAuth.Login({email, password})
            .then(resp => {
                setAuth(resp);

                setCookie(cookieAuthToken, resp.token);
                navigationGo("/");
            })
            .catch((err) => {
                // TODO think about error codes
                if (err === "InvalidCredentials") {
                    setError('Invalid email or password')
                } else {
                    setError('Something wrong, please try again later')
                }
            })
    };

    const onRegister = () => {
        ApiAuth.Register({email, password})
            .then(resp => {
                setAuth(resp);

                setCookie(cookieAuthToken, resp.token);
                navigationGo("/");
            })
            .catch((err) => {
                if (err === "EmailAlreadyRegistered") {
                    setError('This email already registered')
                } else {
                    setError('Something wrong, please try again later')
                }
            })
    }

    if (isVK) {
        return <div>Authorizing via VK...</div>;
    }

    return <div>
        {isReg ?
            <h1>Registration</h1> :
            <h1>Login</h1>
        }

        <div className={styles.formRow}>
            <div className={styles.formLeft}>Email:</div>
            <input type="text" value={email} onChange={e => setEmail(e.target.value)}/>
        </div>

        <div className={styles.formRow}>
            <div className={styles.formLeft}>Password:</div>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}/>
        </div>

        <div className={styles.formRow}>
            <button type="button" onClick={onAuth}>Login</button>
        </div>

        {error &&
            <div className={styles.error}>{error}</div>
        }

        {!isReg &&
            <div className={styles.formRow}>
                <a href={vkAuthURL}>Login via VK</a>
            </div>
        }

        <div className={styles.formRow}>
            {!isReg && <Link href="/auth?registration">Register new account</Link>}
            {isReg && <Link href="/auth">Login to existing account</Link>}
        </div>
    </div>
}
