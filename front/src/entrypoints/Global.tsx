import ReactDOM from 'react-dom';
import React from "react";
import {Root} from "../components/Root/Root";

if (!window.rootLoaded) {
    window.rootLoaded = true;
    ReactDOM.render(<Root/>, document.getElementById('root'));
}
