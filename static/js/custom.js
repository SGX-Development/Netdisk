"use strict";
/*
 Theme Name: Vedha Corporate - Bootstrap HTML template
 Description: Custom JS are defined in this class
 Author: Jyostna
 Author URI: http://themeforest.net/user/jyostna
 Version: 1.0

 -------------------------------------------- */
/*
 TABLE OF CONTENT
 -------------------------------------------------
 1- HEADER
 2- REVOLUTIONARY SLIDER
 2.1- HOME
 2.2- ABOUT US, CAREERS AND INDUSTRIES
 3- STYLES PAGE
 4- I check
 5- Select2
 6- INSIGHTS
 6.1- OWL-CAROUSEL
 7- ABOUT US
 8- SERVICES
 9- WOW
 10- PROGRESS BARS
 11- CAREERS
 12- bLOG
 13- NEWSTICKER
 14- BLOG SINGLE POST
 15- BLOG VIDEO POST
 16- GALLERY AND PORTFOLIO
 17- PRODUCTS
 18- TWITTER FEED
 19- MAPS
 20- CASE STUDIES
 21- BACK TO TOP BUTTON
 22- PRELOADER
 23- Mail Chimp

 -------------------------------------*/

$(window).on('load', function () {
    /*---------Preloader------------*/
    $('.preloader img, .preloader').fadeOut();
    /*------------Isotope Triggering------------*/
    $(".active1").click();
});
$(document).ready(function () {

    /*-------------Header--------------*/
    $('ul.dropdown-menu [data-toggle=dropdown]').on('click', function (event) {
        event.preventDefault();
        event.stopPropagation();
        $(this).parent().siblings().removeClass('open');
        $(this).parent().toggleClass('open');
    });
    $('#myCarousel').carousel({
        interval: 2000
    });
    var windowWidth = $(window).width();
    if (windowWidth > 768) {
        $(".dropdown-toggle").on('mouseenter', function () {
            $(this).click();
        });
        $(".list2").on('mouseenter', function () {
            $(".dropdown").removeClass("open");
        });
        $("header").on('mouseleave', function () {
            $(".dropdown").removeClass("open");

        });
        var dropdownSelectors = $('.dropdown, .dropup');

        dropdownSelectors.on({
            "show.bs.dropdown": function () {
                // On show, start in effect
                var dropdown = dropdownEffectData(this);
                dropdownEffectStart(dropdown, dropdown.effectIn);
            },
            "shown.bs.dropdown": function () {
                // On shown, remove in effect once complete
                var dropdown = dropdownEffectData(this);
                if (dropdown.effectIn && dropdown.effectOut) {
                    dropdownEffectEnd(dropdown, function () {
                    });
                }
            },
            "hide.bs.dropdown": function (e) {
                // On hide, start out effect
                var dropdown = dropdownEffectData(this);
                if (dropdown.effectOut) {
                    e.preventDefault();
                    dropdownEffectStart(dropdown, dropdown.effectOut);
                    dropdownEffectEnd(dropdown, function () {
                        dropdown.dropdown.removeClass('open');
                    });
                }
            }
        });
    }
    function dropdownEffectData(target) {

        var effectInDefault = null,
            effectOutDefault = null;
        var dropdown = $(target),
            dropdownMenu = $('.dropdown-menu', target);
        var parentUl = dropdown.parents('ul.nav');
        if (parentUl.size() > 0) {
            effectInDefault = parentUl.data('dropdown-in') || null;
            effectOutDefault = parentUl.data('dropdown-out') || null;
        }

        return {
            target: target,
            dropdown: dropdown,
            dropdownMenu: dropdownMenu,
            effectIn: dropdownMenu.data('dropdown-in') || effectInDefault,
            effectOut: dropdownMenu.data('dropdown-out') || effectOutDefault
        };
    }

    function dropdownEffectStart(data, effectToStart) {
        if (effectToStart) {
            data.dropdown.addClass('dropdown-animating');
            data.dropdownMenu.addClass('animated');
            data.dropdownMenu.addClass(effectToStart);
        }
    }

    function dropdownEffectEnd(data, callbackFunc) {
        var animationEnd = 'webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend';
        data.dropdown.one(animationEnd, function () {
            data.dropdown.removeClass('dropdown-animating');
            data.dropdownMenu.removeClass('animated');
            data.dropdownMenu.removeClass(data.effectIn);
            data.dropdownMenu.removeClass(data.effectOut);

            if (typeof callbackFunc == 'function') {
                callbackFunc();
            }
        });
    }

    $(".dropdown-submenu>.dropdown-menu").on("mouseenter mouseleave", function () {
        $(".service_nested").click();
    });

    $(".register").on("click", function () {
        $(".signup").addClass("active in");
        $(" .signin").removeClass("active in");
        $("#reg_tab").addClass("active");
        $("#signin_tab").removeClass("active");

    });
    $(".signin").on("click", function () {
        $(".signup").removeClass("active in");
        $(" .signin").addClass("active in");
        $("#reg_tab").removeClass("active");
        $("#signin_tab").addClass("active");

    });
    /*----------header remove hover effect--------------*/
    $(".dropdown").on('mouseover', function () {
        $(this).find(".menu_hover").addClass("hvr-underline-from-center");
    }).on('mouseleave', function () {
        $(this).find(".menu_hover").removeClass("hvr-underline-from-center");
    });


    /*-------------Home Revolutionary slider-------------*/
    if (typeof revslider_showDoubleJqueryError == 'function') {
        $(".home_revolution_slider .tp-banner-slider, .home_slider").revolution({
            dottedOverlay: "none",
            delay: 5000,
            startwidth: 1170,
            startheight: 550,
            hideThumbs: 200,
            hideTimerBar: "on",
            thumbWidth: 100,
            thumbHeight: 50,
            thumbAmount: 5,

            navigationType: "bullet",
            navigationArrows: "none",
            navigationStyle: "preview2",
            touchenabled: "on",
            onHoverStop: "on",
            swipe_velocity: 0.7,
            swipe_min_touches: 1,
            swipe_max_touches: 1,
            drag_block_vertical: false,
            parallax: "mouse",
            parallaxBgFreeze: "on",
            parallaxLevels: [10, 9, 8, 7, 4, 3, 2, 5, 4, 3, 2, 1, 0],
            keyboardNavigation: "on",
            navigationHAlign: "center",
            navigationVAlign: "bottom",
            navigationHOffset: 0,
            navigationVOffset: 20,
            soloArrowLeftHalign: "left",
            soloArrowLeftValign: "center",
            soloArrowLeftHOffset: 20,
            soloArrowLeftVOffset: 0,
            soloArrowRightHalign: "right",
            soloArrowRightValign: "center",
            soloArrowRightHOffset: 20,
            soloArrowRightVOffset: 0,
            shadow: 0,
            fullWidth: "off",
            fullScreen: "off",
            spinner: "spinner4",
            stopLoop: "off",
            stopAfterLoops: -1,
            stopAtSlide: -1,
            shuffle: "off",
            autoHeight: "on",
            forceFullWidth: "on",
            hideThumbsOnMobile: "off",
            hideNavDelayOnMobile: 1500,
            hideBulletsOnMobile: "off",
            hideArrowsOnMobile: "off",
            hideThumbsUnderResolution: 0,
            hideSliderAtLimit: 0,
            hideCaptionAtLimit: 0,
            hideAllCaptionAtLilmit: 0,
            startWithSlide: 0,
            fullScreenOffsetContainer: ".header"
        });
    }
    /*--------- About Us, Careers and Industries Revolutionary slider--------*/
    if (typeof revslider_showDoubleJqueryError == 'function') {
        $(".tp-banner-slider").revolution({
            dottedOverlay: "none",
            delay: 4000,
            startwidth: 1170,
            startheight: 350,
            hideThumbs: 200,
            hideTimerBar: "on",
            thumbWidth: 100,
            thumbHeight: 50,
            thumbAmount: 5,

            navigationType: "bullet",
            navigationArrows: "none",
            navigationStyle: "preview2",

            touchenabled: "on",
            onHoverStop: "on",

            swipe_velocity: 0.7,
            swipe_min_touches: 1,
            swipe_max_touches: 1,
            drag_block_vertical: false,

            parallax: "mouse",
            parallaxBgFreeze: "on",
            parallaxLevels: [7, 4, 3, 2, 5, 4, 3, 2, 1, 0],

            keyboardNavigation: "on",

            navigationHAlign: "center",
            navigationVAlign: "bottom",
            navigationHOffset: 0,
            navigationVOffset: 20,

            soloArrowLeftHalign: "left",
            soloArrowLeftValign: "center",
            soloArrowLeftHOffset: 20,
            soloArrowLeftVOffset: 0,

            soloArrowRightHalign: "right",
            soloArrowRightValign: "center",
            soloArrowRightHOffset: 20,
            soloArrowRightVOffset: 0,

            shadow: 0,
            fullWidth: "off",
            fullScreen: "off",

            spinner: "spinner4",

            stopLoop: "off",
            stopAfterLoops: -1,
            stopAtSlide: -1,

            shuffle: "off",

            autoHeight: "on",
            forceFullWidth: "on",


            hideThumbsOnMobile: "off",
            hideNavDelayOnMobile: 1500,
            hideBulletsOnMobile: "off",
            hideArrowsOnMobile: "off",
            hideThumbsUnderResolution: 0,

            hideSliderAtLimit: 0,
            hideCaptionAtLimit: 0,
            hideAllCaptionAtLilmit: 0,
            startWithSlide: 0,
            fullScreenOffsetContainer: ".header"
        });
    }

    /*----------styles page---------------*/
    $(".panel-collapse").on('shown.bs.collapse', function () {
        $(this).parent(".panel").find(".accordion-section-title>div > .fa").removeClass("fa-plus").addClass("fa-minus");
    }).on('hide.bs.collapse', function () {
        $(this).parent(".panel").find(".accordion-section-title>div > .fa").removeClass("fa-minus").addClass("fa-plus");
    });
    /*--------------icheck ---------------*/
    if ($.fn.iCheck !== undefined) {
        $('.red').iCheck({
            checkboxClass: 'icheckbox_minimal-red',
            radioClass: 'iradio_square-red'
        });
        $('.blue').iCheck({
            checkboxClass: 'icheckbox_minimal-blue'
        });
    }
    /*------------Select2 Plugin------------*/
    if ($.fn.select2 !== undefined) {
        $('.select2').select2();
    }
    /*--------------Insights------------------*/

    /*-------------owl carousel-----------*/
    if ($.fn.owlCarousel !== undefined) {
        $('.insights_owl_carousel .owl-carousel').owlCarousel({
            loop: true,
            margin: 20,
            autoplay: true,
            autoplayTimeout: 2000,
            autoplayHoverPause: true,
            responsive: {
                0: {
                    items: 1
                },
                768: {
                    items: 2
                },
                1024: {
                    items: 3
                }
            }
        }).prepend($('.owl-nav, .owl-dots'));
    }

    /*--------------About-us------------------*/
    if ($.fn.owlCarousel !== undefined) {
        $('.about_carousel .owl-carousel').owlCarousel({
            loop: true,
            margin: 20,
            autoplay: true,
            Height: 300,
            autoplayTimeout: 3000,
            autoplayHoverPause: true,
            responsiveClass: true,
            responsive: {
                0: {
                    items: 1,
                    nav: false,
                    loop: true
                },
                400: {
                    items: 2,
                    nav: false
                },
                500: {
                    items: 3,
                    nav: false
                },
                600: {
                    items: 4,
                    nav: false
                },
                1000: {
                    items: 6,
                    nav: false,
                    loop: true
                }
            }
        }).prepend($('.owl-nav, .owl-dots'));
    }
    $(".profile").on('mouseenter', function () {
        $(this).find("img").addClass("animated pulse");
    }).on('mouseleave', function () {
        $(this).find("img").removeClass("animated pulse");

    });
    /*-------------About Us---------------*/
    if ($(".circliful_custom")[0]) {
        $("#test-circle").circliful({
            animation: 1,
            animationStep: 1,
            foregroundBorderWidth: 8,
            backgroundBorderWidth: 8,
            percent: 87,
            foregroundColor: '#86cb35',
            backgroundColor: '#fff',
            textSize: 15,
            textStyle: 'font-size: 12px;',
            fontColor: '#fff',

            multiPercentage: 1,
            percentages: [10, 20, 30],
            animateInView: true
        });
        $("#test-circle1").circliful({
            animation: 1,
            animationStep: 1,
            foregroundBorderWidth: 8,
            backgroundBorderWidth: 8,
            percent: 95,
            foregroundColor: '#86cb35',
            backgroundColor: '#fff',
            textSize: 15,
            textStyle: 'font-size: 12px;',
            fontColor: '#fff',
            multiPercentage: 1,
            percentages: [10, 20, 30],
            animateInView: true
        });
        $("#test-circle2").circliful({
            animation: 1,
            animationStep: 1,
            foregroundBorderWidth: 8,
            backgroundBorderWidth: 8,
            foregroundColor: '#86cb35',
            backgroundColor: '#fff',
            percent: 70,
            textSize: 15,
            textStyle: 'font-size: 12px;',
            fontColor: '#fff',
            multiPercentage: 1,
            percentages: [10, 20, 30],
            animateInView: true
        });

    }


    /*---------------Services---------------------*/
    if ($.fn.owlCarousel !== undefined) {
        $('.services_carousel .owl-carousel').owlCarousel({
            loop: true,
            margin: 20,
            autoplay: true,
            Height: 300,
            autoplayTimeout: 3000,
            autoplayHoverPause: true,
            responsiveClass: true,
            responsive: {
                0: {
                    items: 1,
                    nav: false,
                    loop: true
                },
                400: {
                    items: 2,
                    nav: false
                },
                600: {
                    items: 3,
                    nav: false
                },
                992: {
                    items: 4,
                    nav: false,
                    loop: true
                }
            }
        });
        $('.services_carousel.owl-carousel').prepend($('.owl-nav, .owl-dots'));
    }


    /*------------------Wow initialisation----------*/
    if (typeof WOW === "function") {
        new WOW().init();
    }
    /*--------------Progress bars-------------------*/
    $(".progress-bar-success").animate({
        width: "92%"
    }, 4000);
    $(".progress-bar-primary").animate({
        width: "70%"
    }, 4000);
    $(".progress-bar-warning").animate({
        width: "30%"
    }, 4000);
    $(".progress-bar-danger").animate({
        width: "15%"
    }, 4000);
    $(".progress-bar-info").animate({
        width: "75%"
    }, 4000);

    /*-----------Blog page----------------*/
    var currentValue;
    $("#bloglike1, #bloglike2, #bloglike3,#bloglike_bs").one("click", function () {
        currentValue = $(this).find("span").text();

    }).on("click", function () {
        var val1 = ($(this).find("span").text());

        if (val1 == currentValue) {
            val1++;
            $(this).find("span").text(val1);

        }
        else {
            val1--;
            $(this).find("span").text(val1);
        }
        $(this).toggleClass("text-primary");
    });

    /*-------------NewsTicker---------------*/
    if ($.fn.newsTicker !== undefined) {

        $('.newsticker').newsTicker({
            row_height: 75,
            max_rows: 3,
            duration: 1500
        });
        /*-----------------Home-----------------*/
        $('.newssticker_home').newsTicker({
            row_height: 165,
            max_rows: 3,
            duration: 2000
        });

    }

    /*------------Hover fucnctionality---------------*/
    $('.blog_links').on('mouseenter', function () {
        $(this).addClass('blog_margin');
    }).on('mouseleave', function () {
        $(this).removeClass('blog_margin');
    });

    $('.blog_links_sd').on('mouseenter', function () {
        $(this).addClass('blog_margin_sd');
    }).on('mouseleave', function () {
        $(this).removeClass('blog_margin_sd');
    });
    /*---------------Blog single----------------*/


    /*-----------Blog Video Post----------------*/

    if ($.fn.owlCarousel !== undefined) {
        $('.blogvideo_carousel .owl-carousel').owlCarousel({
            loop: true,
            nav: false,
            margin: 20,
            Height: 300,
            autoplayTimeout: 3000,
            autoplayHoverPause: true,
            responsiveClass: true,
            responsive: {
                0: {
                    items: 1

                },
                600: {
                    items: 2

                },
                992: {
                    items: 3

                }
            }


        }).prepend($('.owl-nav'));
        $('.owl-stage-outer').prepend($('.owl-dots'));
        /*---------------------Index 2----------------*/
        $('.owl-carousel').owlCarousel({
            loop: true,
            margin: 20,
            autoplay: true,
            autoplayTimeout: 2000,
            autoplayHoverPause: true,
            navs: true,
            responsive: {
                0: {
                    items: 1
                },
                768: {
                    items: 2
                },
                1024: {
                    items: 3
                }
            }
        });
    }
    /*------------Gallery & portfolio----------------*/
    if ($.fn.isotope !== undefined) {
        var $container = $('#gallery, #portfolio, #posts').isotope({
            itemSelector: '.isotope-item',
            isFitWidth: true
        });

        $container.isotope({
            filter: '*'
        });
        $('#filters').on('click', 'li', function () {
            var filterValue = $(this).attr('data-filter');
            $container.isotope({
                filter: filterValue
            });
            $(this).closest("ul").find(".active").removeClass("active");
            $(this).find("a").addClass("active");
        });
    }

    if ($.fn.lightbox !== undefined) {
        lightbox.option({
            'resizeDuration': 500,
            'wrapAround': true
        });
    }
    /*-------------Products--------------*/

    $(".icon_show").on("click", function () {
        $(this).hide();
        $(".second_section").show();

    });

    /*------------Twitter feed------------------*/

    $('.tweet').twittie({
        dateFormat: '%b. %d, %Y',
        template: '<i class="fa fa-twitter" aria-hidden="true"></i>&nbsp;{{tweet}} <div class="date">{{date}}</div><br />',
        count: 2,
        hideReplies: true,
        apiPath: 'twitter_api/tweet.php'
    });

    /*--------------maps--------------------------*/
    if ($.fn.gmap3 !== undefined) {
        $(".home_map").gmap3({
            map: {
                options: {
                    center: [17.400408, 78.507905],
                    zoom: 8,
                    scrollWheel: false
                }
            },
            marker: {
                values: [{
                    address: "nallakunta,hyderabad ",
                    options: {
                        icon: "image/location-mark.png"
                    }
                }]
            }
        });
    }
    if ($.fn.gmap3 !== undefined) {
        $(".map").gmap3({
            map: {
                options: {
                    center: [17.400408, 78.507905],
                    zoom: 8,
                    scrollwheel: false,
                    styles: [
                        {
                            "featureType": "all",
                            "stylers": [
                                {
                                    "saturation": 0
                                },
                                {
                                    "hue": "#6fbd11"
                                }
                            ]
                        },
                        {
                            "featureType": "road",
                            "stylers": [
                                {
                                    "saturation": -70
                                }
                            ]
                        },
                        {
                            "featureType": "transit",
                            "stylers": [
                                {
                                    "visibility": "off"
                                }
                            ]
                        },
                        {
                            "featureType": "poi",
                            "stylers": [
                                {
                                    "visibility": "off"
                                }
                            ]
                        },
                        {
                            "featureType": "water",
                            "stylers": [
                                {
                                    "visibility": "simplified"
                                },
                                {
                                    "saturation": -60
                                }
                            ]
                        }
                    ]


                }
            },
            marker: {
                values: [{
                    address: "nallakunta,hyderabad ",
                    options: {
                        icon: "image/location-mark.png"
                    }
                }]
            }
        });
    }
    /*------------case studies hover effect-------------------*/
    $(".case-item").on('mouseenter', function () {
        $(this).find(".text-justify").addClass("animated fadeOut");
        $(this).find(".btn-empty").css("top", "-50px").addClass("animated fadeInUp");
        $(this).find(".case_image").css("top", "70px").addClass("animated fadeInDown").removeClass("animated fadeInUp");
        $(this).find("img").addClass("animated fadeInDown");
        $(this).find("h5").addClass("animated fadeInDown");
    }).on('mouseleave', function () {
        $(this).find(".text-justify").removeClass("animated fadeOut").addClass("animated pulse");
        $(this).find(".btn-empty").css("top", "0").removeClass("animated fadeInUp").addClass("animated fadeInDown");
        $(this).find(".case_image").css("top", "0").removeClass("animated fadeInDown").addClass("animated fadeInUp");
        $(this).find("img").removeClass("animated fadeInDown");
        $(this).find("h5").removeClass("animated fadeInDown");
    });
    /*-------------Back to top button------------------*/
    $('#return-to-top').fadeOut().on('click', function () {
        $('body,html').animate({
            scrollTop: 0
        }, 500);
    });
    $(window).on('scroll', function () {
        if ($(this).scrollTop() >= 50) {
            $('#return-to-top').fadeIn(200);
        } else {
            $('#return-to-top').fadeOut(200);
        }
    });
    /*------------------404 page ---------------------*/
    if ($(window).width() > 991) {
        $('.body_back').mousemove(function (e) {
            var x = -(e.pageX + this.offsetLeft) / 20;
            var y = -(e.pageY + this.offsetTop) / 20;
            $(this).css('background-position', x + 'px ' + y + 'px');
        });
    }
    setTimeout(function () {
        $(".rotated").addClass("rotate");
    }, 3000);

    /*--------------mail chimp----------------*/
    $('#subscribe').on('submit', function () {

        if (!valid_email_address($("#email").val())) {
            swal({
                type: 'error',
                html: 'Please Enter Valid Email Address'
            })
        }
        else {
            $.ajax({
                url: 'subscribe.php',
                data: $('#email').serialize(),
                type: 'POST',
                success: function () {
                    swal({
                        type: 'success',
                        html: 'You have Successfully Subscribed'
                    });
                }
            });
        }

        return false;
        function valid_email_address(email) {
            var pattern = new RegExp(/^[+a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/i);
            return pattern.test(email);
        }

    });
    //=============Swiper==============
    if (typeof Swiper == 'function') {
        /*------------Index 2-------------------*/
        var swiper5 = new Swiper('.products_slider', {
            pagination: '.swiper-pagination',
            slidesPerView: 6,
            paginationClickable: true,
            loop: true,
            freeMode: true,
            autoplay: 2000,
            breakpoints: {

                320: {
                    slidesPerView: 1,
                    spaceBetweenSlides: 10
                },

                480: {
                    slidesPerView: 2,
                    spaceBetweenSlides: 20
                },

                640: {
                    slidesPerView: 3,
                    spaceBetweenSlides: 30
                }
            }
        });
        var swiper = new Swiper('.media-swiper .swiper-container', {
            pagination: '.swiper-pagination',
            paginationClickable: true,
            spaceBetween: 30,
            autoplay: 2500,
            centeredSlides: true,
            autoplayDisableOnInteraction: false
        });
        var swiper35 = new Swiper('.Products_page_swiper', {
            effect: 'flip',
            grabCursor: true,
            loop: true,
            autoplay: 3000
        });
    }
    /*------------------------Dribbble widget----------------------*/
    if ($(".shots")[0]) {
        $.jribbble.setToken('f688ac519289f19ce5cebc1383c15ad5c02bd58205cd83c86cbb0ce09170c1b4');
        $.jribbble.shots('debuts', {'per_page': 6, 'timeframe': 'month', 'sort': 'views'}).then(function (res) {
            var html = [];
            res.forEach(function (shot) {
                html.push('<li class="shots--shot">');
                html.push('<a href="' + shot.html_url + '" target="_blank">');
                html.push('<img src="' + shot.images.normal + '">');
                html.push('</a></li>');
            });
            $('.shots').html(html.join(''));
        });
    }
    /*-------------plyr----------------------*/
    if ($('.video')[0]) {
        plyr.setup();
    }
    /*-------------Count Up-------------*/
    if (typeof CountUp === "function") {
        var options = {
            useEasing: true,
            useGrouping: true,
            separator: '',
            decimal: '.',
            prefix: '',
            suffix: ''
        };
        /*-------------Index---------------*/
        if ($("#Overseas")[0]) {
            var demo4 = new CountUp("Success", 0, 97, 0, 7, options);
            var demo3 = new CountUp("Client", 0, 1400, 0, 7, options);
            var demo2 = new CountUp("Investors", 0, 1600, 0, 7, options);
            var demo1 = new CountUp("Overseas", 0, 56, 0, 7, options);
            $(window).on('scroll', function () {
                var winTop = $(window).scrollTop();
                var winHeight = $(window).height();

                var animation1 = $('#Overseas').offset().top;
                if (winTop >= (animation1 - winHeight)) {
                    demo4.start();
                    demo3.start();
                    demo2.start();
                    demo1.start();


                }
            });
        }
        /*------------Carousel Time interval-------------*/
        $('#services_carousel').carousel({
            interval: 3000
        });

    }

});







(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
    (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

ga('create', 'UA-59850948-1', 'auto');
ga('send', 'pageview');
