/**
 * @file scss tool.
 *
 * @author <a href='http://vanessa.b3log.org'>Liyuan Li</a>
 * @version 0.2.0.0, Oct 3, 2018
 */

const gulp = require('gulp')
const sass = require('gulp-sass')

function sassProcess () {
  return gulp.src('./scss/*.scss').
    pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError)).
    pipe(gulp.dest('./css'))
}

function sassProcessWatch () {
  gulp.watch('./scss/*.scss', sassProcess)
}

gulp.task('watch', gulp.series(sassProcessWatch))