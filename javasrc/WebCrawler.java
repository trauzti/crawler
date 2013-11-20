// A simple Web Crawler written in Java
// Partly based on http://cs.nyu.edu/courses/fall02/G22.3033-008/proj1.html

// Our crawler is a sequential, topical crawler.
// It uses a priority queue for unvisited links.
// The links are scored with 1.0 if the topic is part of the anchor text for a link, else 0.0.
// This crawler lists HTML pages that contain the query words in their <body> section.
// Usage: From command line 
//     java WebCrawler <URL> <TOPIC> <QUERY WORDS> <N>
//  where   URL is the url (seed) to start the crawl,
//          TOPIC is the topic we are interested in (used to quide the crawler to relevant links)
//          QUERY WORDS is the phrase we are interested in
//          N (optional) is the maximum number of pages to crawl
//  Hrafn Loftsson, Reykjavik University, Fall 2013
//import java.text.*;

import java.util.*;
import java.io.*;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;

public class WebCrawler {
    private final int SEARCH_LIMIT = 500;  	// Absolute max # of pages crawled. Respect this, be polite!
    private final int MILLISECOND_WAIT = 300;   // Wait between url requests, be polite!
    private final boolean DEBUG = false;        // To control debugging output
    private final String userAgent = "RuBot"; 	// Reykjavik University bot

    Frontier frontier;      // The frontier, the list of pages yet to be crawled (visited)
    Hashtable<String, Integer> visitedURLs;    // The list of visited URLs
    URLCanonicalizer canonicalizer; // Used to transform URLs to canonical form
    int maxPages;           // max number of pages to crawl, may be supplied by the user
    String topic;           // the topic we are interested in
    String queryString;    // the query string we are interested in
    String[] queryWords;   // individual words of the query string

    RobotTxtParser robotParser; // A robots.txt parser
    HTMLParser htmlParser;  	// A HTMLParser
    int totalRelevant=0;    	// Total number of pages containing our query string



    public void initialize(String[] argv) {
        String url;
        robotParser = new RobotTxtParser(userAgent, DEBUG);
        htmlParser = new HTMLParser();
        canonicalizer= new URLCanonicalizer();

        visitedURLs = new Hashtable<String, Integer>();
        frontier = new Frontier(DEBUG);

        url = argv[0];                  // The seed URL supplied by the user
        topic = argv[1].toLowerCase();  // The topic
        queryString = argv[2].toLowerCase().replaceAll("\\s+", " ");	// The query words supplied by the user
        queryWords = queryString.split("\\s");    // Assume space between query words

        String canonicalUrl = canonicalizer.getCanonicalURL(url);	// Canonicalize the URL
        frontier.add(canonicalUrl, 0.0);                            // The seed has score 0.0

        maxPages = SEARCH_LIMIT;
        if (argv.length > 3) { // Does the user override the search limit?
            int iPages = Integer.parseInt(argv[3]);
            if (iPages < maxPages)
                maxPages = iPages;
        }

        System.out.println("--------------------------------------------------------");
        System.out.println("Starting crawl, seed: " + canonicalUrl);
        System.out.println("Topic: " + topic);
        System.out.println("Query string: " + queryString);
        System.out.println("Maximum number of pages to visit: " + maxPages);
        System.out.println("--------------------------------------------------------");
   }   

    // Retrieve the links (href) from the given url
    private Elements getLinks(String url)
    {
        Elements links;
        try {
            links = htmlParser.getLinks();      // Retrieve the <a href> links
        }
        catch (IOException e) {
            System.out.println("Could not get links from " + url);
            links = new Elements();             // Empty elements
        }
        return links;
    }

    // Adds the retrieved links to the frontier
    private void addLinks(Elements links)
    {
	/********************************************************/
	/* GAP!							*/
	/* Make sure that you add canonicalized versions	*/
	/* of the links to the frontier				*/
	/* You also need to score the links			*/
	/********************************************************/
    }

	// Returns true if our phrase query is found in the given text, otherwise false.
    private boolean isRelevantText(String text) {
	/********************************************************/
	/* GAP!							*/
	/* You do not have to implement stemming.		*/
	/* However, make the comparison case-insensitive.	*/
	/********************************************************/
    }
    private boolean isRelevantUrl(String url) {
        /********************************************************/
	/* GAP!							*/
	/* Returns true if the body of the page			*/
	/* corresponding to the url is relevant, 		*/
	/* i.e. if it contains the phrase query			*/
	/* Uses the relevantText() method			*/
	/********************************************************/
    }

    private void processUrl(String url)
    {
	/********************************************************/
	/* GAP!							*/
	/* Process the given url, which means at least:		*/
	/* 1) Connect to it using the HTML parser		*/
	/* 2) Print an appropriate message if it is relevant	*/
	/* 3) Extract links from the url and add to frontier	*/
	/********************************************************/	
    }

    private void wait(int milliseconds)	// Halt execution for the specified number of milliseconds
    {
        try {
            Thread.currentThread().sleep(milliseconds);
        }
        catch (InterruptedException e) {
                e.printStackTrace();
        }
    }

    // This method does the crawling.
    public void crawl()
    {
        for (int i = 0; i < maxPages; i++) {        // Visit maxPages

	    /****************************************************/
	    /* GAP!						*/
	    /* Retrive the next url from the frontier		*/
	    /* parse and process it given that			*/
	    /* 	you are allowed to do so (robot.txt)		*/
	    /* 	and that the url has not been visitied before 	*/
	    /****************************************************/

            wait(MILLISECOND_WAIT);          // Be polite, wait x milliseconds before fetching the next page

            if (frontier.isEmpty()) break;   // Break the loop if frontier is empty
        }
        System.out.println("--------------------------------------------------------");
        System.out.println("Search complete, " + maxPages + " pages crawled");
        System.out.println("Search query " + queryString + " found in " + totalRelevant + " pages");
        System.out.println("Total distinctive urls found: " + frontier.totalCount());
        System.out.println("--------------------------------------------------------");
    }

    public static void main(String[] argv)
    {
        WebCrawler wc = new WebCrawler();
        wc.initialize(argv);
        wc.crawl();
    }



}

