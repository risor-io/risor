;;; monkey.el --- mode for editing monkey scripts

;; Copyright (C) 2018 Steve Kemp

;; Author: Steve Kemp <steve@steve.fi>
;; Keywords: languages
;; Version: 1.0

;;; Commentary:

;; Provides support for editing monkey scripts with full support for
;; font-locking, but no special keybindings, or indentation handling.

;;;; Enabling:

;; Add the following to your .emacs file

;; (require 'monkey)
;; (setq auto-mode-alist (append '(("\\.mon$" . monkey-mode)) auto-mode-alist)))



;;; Code:

(defvar monkey-constants
  '("true"
    "false"))

(defvar monkey-keywords
  '(
    "else"
    "fn"
    "for"
    "foreach"
    "function"
    "if"
    "in"
    "let"
    "return"
    ))

;; The language-core and functions from the standard-library.
(defvar monkey-functions
  '(
    "args"
    "exit"
    "file.close"
    "file.lines"
    "file.open"
    "first"
    "int"
    "last"
    "len"
    "math.abs"
    "math.random"
    "math.sqrt"
    "push"
    "puts"
    "read"
    "rest"
    "set"
    "string"
    "string.interpolate"
    "string.reverse"
    "string.split"
    "string.tolower"
    "string.toupper"
    "string.trim"
    "type"
    "version"
    ))


(defvar monkey-font-lock-defaults
  `((
     ("\"\\.\\*\\?" . font-lock-string-face)
     (";\\|,\\|=" . font-lock-keyword-face)
     ( ,(regexp-opt monkey-keywords 'words) . font-lock-builtin-face)
     ( ,(regexp-opt monkey-constants 'words) . font-lock-constant-face)
     ( ,(regexp-opt monkey-functions 'words) . font-lock-function-name-face)
     )))

(define-derived-mode monkey-mode fundamental-mode "monkey script"
  "monkey-mode is a major mode for editing monkey scripts"
  (setq font-lock-defaults monkey-font-lock-defaults)

  ;; Comment handler for single & multi-line modes
  (modify-syntax-entry ?\/ ". 124b" monkey-mode-syntax-table)
  (modify-syntax-entry ?\* ". 23n" monkey-mode-syntax-table)

  ;; Comment ender for single-line comments.
  (modify-syntax-entry ?\n "> b" monkey-mode-syntax-table)
  (modify-syntax-entry ?\r "> b" monkey-mode-syntax-table)
  )

(provide 'monkey)
