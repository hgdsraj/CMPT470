$( document ).ready(function() {
    console.log( "ready!" );
    let gradient = $('#gradient');
    console.log(gradient.siblings());
    gradient.attr('href', 'https://tutorialzine.com/2013/10/12-awesome-css3-features-you-can-finally-use');
    const button = $('<input/>').attr({ type: 'button', name:'btn1', value:'a button' });
    $('body').append(button);
    // Hide all paragraphs using a slide up animation over 0.8 seconds
    $( "p" ).slideUp( 7052 );

    rotateableText = $('#rotate');
    button.click(function () {
        rotateableText.toggleClass('h2Overwrite')
    });

});

$( window ).on( "load", function() {
    console.log( "window loaded" );
});
