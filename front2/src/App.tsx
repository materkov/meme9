import ReactDOM from 'react-dom';
import React from "react";
import {Feed} from "./components/Feed/Feed";
import styles from "./App.module.css";
import {Profile} from "./components/Profile/Profile";
import * as schema from "./api/api2";
import {Router} from "./components/Router/Router";

function App() {
    return (
        <div className={styles.Background}>
            <div className={styles.App}>
                <Router/>
            </div>
        </div>
    )
}

window.modules[schema.Renderers.FEED] = Feed
window.modules[schema.Renderers.PROFILE] = Profile

ReactDOM.render(<App/>, document.getElementById('root'));
