/**
 * @file scss tool.
 *
 * @author <a href='http://vanessa.b3log.org'>Liyuan Li</a>
 * @version 0.2.0.0, Oct 3, 2018
 */

const gulp = require('gulp')
const sass = require('gulp-sass')
const concat = require('gulp-concat')
const uglify = require('gulp-uglify')
const cleanCSS = require('gulp-clean-css')

function sassProcess () {
  return gulp.src('./scss/*.scss').
    pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError)).
    pipe(gulp.dest('./css'))
}

function sassProcessWatch () {
  gulp.watch('./scss/*.scss', sassProcess)
}

gulp.task('watch', gulp.series(sassProcessWatch))

function minCodemirrorJS () {
  var jsCodemirror = [
    './js/lib/codemirror-5.40.2/lib/codemirror.js',
    './js/lib/codemirror-5.40.2/addon/edit/closetag.js',
    './js/lib/codemirror-5.40.2/addon/fold/xml-fold.js',
    './js/lib/codemirror-5.40.2/addon/fold/foldcode.js',
    './js/lib/codemirror-5.40.2/addon/fold/foldgutter.js',
    './js/lib/codemirror-5.40.2/addon/fold/brace-fold.js',
    './js/lib/codemirror-5.40.2/addon/fold/xml-fold.js',
    './js/lib/codemirror-5.40.2/addon/fold/indent-fold.js',
    './js/lib/codemirror-5.40.2/addon/fold/comment-fold.js',
    './js/lib/codemirror-5.40.2/mode/xml/xml.js',
    './js/lib/codemirror-5.40.2/mode/javascript/javascript.js',
    './js/lib/codemirror-5.40.2/mode/css/css.js',
    './js/lib/codemirror-5.40.2/mode/htmlmixed/htmlmixed.js',
    './js/lib/codemirror-5.40.2/addon/merge/merge.js',
    './js/lib/js-beautify-1.8.5/beautify.min.js',
    './js/lib/js-beautify-1.8.5/beautify-html.min.js',
    './js/lib/diff_match_patch.js',
  ]
  return gulp.src(jsCodemirror).
    pipe(uglify()).
    pipe(concat('codemirror.min.js')).
    pipe(gulp.dest('./js/lib/codemirror-5.40.2'))
}

function minCodemirrorCSS () {
  return gulp.src([
    './js/lib/codemirror-5.40.2/lib/codemirror.css',
    './js/lib/codemirror-5.40.2/addon/fold/foldgutter.css',
    './js/lib/codemirror-5.40.2/addon/merge/merge.css']).
    pipe(cleanCSS()).
    pipe(concat('codemirror.min.css')).
    pipe(gulp.dest('./js/lib/codemirror-5.40.2'))
}

gulp.task('default', gulp.series(sassProcess, gulp.parallel(minCodemirrorCSS, minCodemirrorJS)))