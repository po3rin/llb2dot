let app = new Vue({
    el: '#app',
    data: {
        dockerfile: 'FROM golang:1.12',
        dots: `strict digraph llb {
            // Node definitions.
            "FROM golang:1.12";
            "sha256:a14...";

            // Edge definitions.
            "FROM golang:1.12" -> "sha256:a14...";
            }`
    },
    computed: {
        image: function(){
            var image = Viz(this.dots, {format: 'svg'})
            return image
        }
    },
    methods: {
        update: function(){
            axios.post('http://localhost:8080/api/dot', this.dockerfile).then(res => {
                this.dots = res.data;
            },
            error => {
                console.error(this.dockerfile)
                console.error(error)
            })
        }
    }
})