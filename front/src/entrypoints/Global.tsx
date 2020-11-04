import ReactDOM from 'react-dom';
import React from "react";
import {Root} from "../components/Root/Root";
import {Error} from "../components/Error/Error";

if (!window.rootLoaded) {
    window.rootLoaded = true;

    window.modules.ErrorRenderer = Error;

    ReactDOM.render(<Root/>, document.getElementById('root'));
}
