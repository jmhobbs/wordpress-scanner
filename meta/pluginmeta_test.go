package meta

import (
	"strings"
	"testing"
)

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsingMetaLinesMatcher(t *testing.T) {
	src := []string{
		// Best case
		" * Plugin Name: example",
		// Lower case
		" * plugin name: example",
		// Weird spacing
		"* 	Plugin Name:    example    ",
	}

	for _, s := range src {
		pm := New()
		if !pm.ParseMetaLine(s) {
			t.Error("Didn't match valid meta line.")
		}
		if pm.Get("Plugin Name") != "example" {
			t.Errorf("Wrong Plugin Name: %s", pm.Get("Plugin Name"))
		}
	}
}

func TestGetSetCasing(t *testing.T) {
	pm := New()
	pm.Set("ALL CAPS", "value")
	if pm.Get("ALL CAPS") != "value" {
		t.Error("Failed to Get 'ALL CAPS'")
	}
	if pm.Get("all caps") != "value" {
		t.Error("Failed to Get 'all caps'")
	}
	if pm.Get("All Caps") != "value" {
		t.Error("Failed to Get 'All Caps'")
	}
}

func TestScan(t *testing.T) {
	r := strings.NewReader(samplePHP)

	pm := New()
	err := pm.Scan(r)
	if err != nil {
		t.Errorf("failed to scan meta: %v", err)
	}

	values := map[string]string{
		"Plugin Name": "bbPress",
		"Plugin URI":  "http://bbpress.org",
		"Description": "bbPress is forum software with a twist from the creators of WordPress.",
		"Author":      "The bbPress Community",
		"Author URI":  "http://bbpress.org",
		"Version":     "2.3-beta1",
		"Text Domain": "bbpress",
		"Domain Path": "/languages/",
	}

	for field, value := range values {
		if pm.Get(field) != value {
			t.Errorf("Invalid value for '%s': %v", field, pm.Get(field))
		}
	}
}

var samplePHP = `<?php

/**
 * The bbPress Plugin
 *
 * bbPress is forum software with a twist from the creators of WordPress.
 *
 * $Id: bbpress.php 4732 2013-01-28 18:30:40Z johnjamesjacoby $
 *
 * @package bbPress
 * @subpackage Main
 */

/**
 * Plugin Name: bbPress
 * Plugin URI:  http://bbpress.org
 * Description: bbPress is forum software with a twist from the creators of WordPress.
 * Author:      The bbPress Community
 * Author URI:  http://bbpress.org
 * Version:     2.3-beta1
 * Text Domain: bbpress
 * Domain Path: /languages/
 */

// Exit if accessed directly
if ( !defined( 'ABSPATH' ) ) exit;

if ( !class_exists( 'bbPress' ) ) :
/**
 * Main bbPress Class
 *
 * "How doth the little busy bee, improve each shining hour..."
 *
 * @since bbPress (r2464)
 */
final class bbPress {

	/** Magic *****************************************************************/

	/**
	 * bbPress uses many variables, several of which can be filtered to
	 * customize the way it operates. Most of these variables are stored in a
	 * private array that gets updated with the help of PHP magic methods.
	 *
	 * This is a precautionary measure, to avoid potential errors produced by
	 * unanticipated direct manipulation of bbPress's run-time data.
	 *
	 * @see bbPress::setup_globals()
	 * @var array
	 */
	private $data;
`
