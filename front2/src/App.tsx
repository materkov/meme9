import ReactDOM from 'react-dom';
import React from "react";
import styles from "./App.module.css";
import {Router} from "./components/Router/Router";
import bridge from "@vkontakte/vk-bridge";

function App() {
    return (
        <div className={styles.Background}>
            <div className={styles.App}>
                <Router/>
            </div>
        </div>
    )
}

ReactDOM.render(<App/>, document.getElementById('root'));

const urlParams = new URLSearchParams(window.location.search);
if (urlParams.get('vk_user_id')) {
    bridge.send("VKWebAppInit", {});
}
