'use strict';

const e = React.createElement;

class LikeButton extends React.Component {
    constructor(props) {
        super(props);
        this.state = { liked: false };
    }

    render() {
        if (this.state.liked) {
            return e(
                'button',
                { onClick: () => this.setState({ liked: false }) },
                'Using react is hard'
            )
        }

        return e(
            'button',
            { onClick: () => this.setState({ liked: true }) },
            'Using react is easy'
        );
    }
}

const domContainer = document.querySelector('#like_button_container');
ReactDOM.render(e(LikeButton), domContainer);