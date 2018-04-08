/**
 * @file scss tool.
 *
 * @author <a href='http://vanessa.b3log.org'>Liyuan Li</a>
 * @version 0.1.0.0, Apr 8, 2018
 */

const gulp = require('gulp')
const sass = require('gulp-sass')

gulp.task('sass', function () {
  return gulp.src('./scss/*.scss')
    .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
    .pipe(gulp.dest('./scss'))
})

gulp.task('watch', function () {
  gulp.watch(['./scss/*.scss'], ['sass'])
})

gulp.task('default', ['sass'])