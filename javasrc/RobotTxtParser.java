/*
robots.txt is a simple file, placed on domain root, that contains directives for robots.
Syntax of robot.txt is simple. Here is an example:

#name of the bad bot
User-agent: BadBot
#disallow all directive for this bot:
Disallow: /

#name of the good bot (not a real name just for example)
User-agent: AdSenseBot
#allow all directive
Disallow:

#rest of the bots
User-agent: *
Disallow: /usr
Disallow: /tmp/
Allow: /tmp/images/
*/

/*
A web crawler calls the parse() method for a given url.  The host is extracted from the url and robots.txt for the host
is parsed and saved in a hash table.  The robots.txt file is thus only parsed once for each host.
 */

import java.io.*;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Hashtable;

public class RobotTxtParser {
    private final String robotTxtString = "robots.txt";
    private final String httpString = "http://";
    private boolean debug = false;
    private String userAgent;       // The parser is interested in info regarding this user agent

    // We hash the directives to make sure that we only parse robots.txt once for each host encountered
    Hashtable<String, RobotDirectives> hostsDirectives;

    public class RobotDirectives {  // Nested class, only used by RobotTxtParser
        ArrayList<String> disallowDirectives;
        ArrayList<String> allowDirectives;

        public RobotDirectives() {
            disallowDirectives = new ArrayList<String>();
            allowDirectives  = new ArrayList<String>();
        }

        public void addDisallow(String directive) {
            disallowDirectives.add(directive);
        }

        public void addAllow(String directive) {
            allowDirectives.add(directive);
        }

        public void clear() {
            disallowDirectives.clear();
            allowDirectives.clear();
        }
    }

    public RobotTxtParser(String userAgent, boolean debug)
    {
        this.userAgent = userAgent;
        this.debug = debug;
        hostsDirectives = new Hashtable<String, RobotDirectives>(50);
    }

    public boolean urlExists(URL url) {

        try
        {
            HttpURLConnection.setFollowRedirects(false);
            HttpURLConnection con = (HttpURLConnection) url.openConnection();
            con.setRequestMethod("HEAD");
            return (con.getResponseCode() == HttpURLConnection.HTTP_OK);
        }
        catch (IOException e) {
            if (debug) System.out.println("RobotTxtParser-Could not make a head request to " + url.toString());
            return false;
        }
  }

    //    Reads in the contents of the stream and returns it as a string
    private String readUrlContent(URL url) {
        InputStream is;
        String data="";

        try {
            is = url.openStream();         // throws an IOException
            BufferedReader d = new BufferedReader(new InputStreamReader(is));

            String s;
            while ((s = d.readLine()) != null)
                    data = data + s + "\n";

            is.close();
         }
        catch (IOException e) {
            if (debug) System.out.println("RobotTxtParser-Could not read " + url.toString());
        }
        return data;
    }


    private String getHost(String url)
    {
        URL urlObject;
        try {
            urlObject = new URL(url);
        }
        catch (MalformedURLException e) {
            if (debug) System.out.println("RobotTxtParser-Malformed URL: " + url);
            return "";
        }

        String host = urlObject.getHost();
        if (host.startsWith("www."))
            host = host.replaceFirst("www.", "");
        return host;
    }

    /*
        Connects to a robot.txt file on a web server and returns its content as a string
    */
    private String getContent(String host)
    {
	    // form URL of the robots.txt file
        String rootPage = httpString + host;
        String robotPage = rootPage + "/" + robotTxtString;
        URL urlRobot;
        String robotTxtData="";

        try {
            urlRobot = new URL(robotPage);
            // Check if the url exists, so we don't get a HTTP 404 response (file not found)
            if (urlExists(urlRobot))
                    robotTxtData = readUrlContent(urlRobot);
            else
                if (debug) System.out.println("RobotTxtParser-robots.txt not found at " + robotPage);

        } catch (MalformedURLException e) {
                  if (debug) System.out.println("RobotTxtParser-Malformed URL: " + robotPage);
        }

        if (debug && !robotTxtData.equals(""))
        {
            System.out.println("Robots.txt at " + host + ":");
            System.out.println("-----------------------------------------");
            System.out.println(robotTxtData);
            System.out.println("-----------------------------------------");
        }

        return robotTxtData;
    }


    /*
    This parser assumes that the contents of the robots.txt file is contained in the argument
    Returns RobotDirectives or null if no directives found for us
     */
    private RobotDirectives parseContent(String robotTxtContent) {

        RobotDirectives robotDirectives = new RobotDirectives();
        if (robotTxtContent.equals(""))     // empty content
            return robotDirectives;

        //split
        String[] lines = robotTxtContent.split("\n");

        int i=0;
        boolean last=false;
        while (i < lines.length && !last) {
            if (lines[i].startsWith("User-agent:")) {
                String agent = lines[i];
                agent = agent.replace("User-agent:", "");    // Keep only the relevant info
                agent = agent.trim();

                // Does this apply to us?
                if (agent.equalsIgnoreCase(userAgent) || agent.equals("*")) {

                    for (int j=i+1; j<lines.length && !last; j++) {
                        if (lines[j].startsWith("Disallow:")) {
                            String directive = lines[j].trim();
                            directive = directive.replace("Disallow:", "");    // Keep only the relevant info
                            directive = directive.trim();
                            robotDirectives.addDisallow(directive);

                        }
                        else if (lines[j].startsWith("Allow:")) {
                            String directive = lines[j].trim();
                            directive = directive.replace("Allow:", "");    // Keep only the relevant info
                            directive = directive.trim();
                            robotDirectives.addAllow(directive);
                        }
                        else if (lines[j].startsWith("User-agent:"))    // Then we have found the start of another agent
                            last = true;                                // We are not interested in that record

                        if (j==lines.length-1)  // Last line?
                                last=true;
                    }
                }
            }
            i++;
        }
        return robotDirectives;
    }

    /* A crawler calls this method to parse the content of robots.txt at the host owning the url */
    public void parse(String url) {

        String host = getHost(url);
        if (!hostsDirectives.containsKey(host))     // If host does not exists in the hash, then parse robots.txt
        {
            String content = getContent(host);
            hostsDirectives.put(host, parseContent(content));   // parse the content and save it into the hash
        }
    }


    public boolean isUrlAllowed(String url) {

        RobotDirectives directives=null;
        String host = getHost(url);

        if (hostsDirectives.containsKey(host))     // If host does not exists in the hash, then parse robots.txt
            directives = hostsDirectives.get(host);
        else
            if (debug) System.out.println("RobotTxtParser-No directives found for host " + host);   // Should not happen!

        if (directives != null) {
            // Is the url specifically allowed?
            for (String directive : directives.allowDirectives ) {
                if (url.contains(directive))
                    return true;
            }

            // Is the url specifically disallowed?
            for (String directive : directives.disallowDirectives ) {
                if (url.contains(directive))
                    return false;
            }
        }

        return true;    // If we get here then no directives were found or the url is neither specifically allowed nor specifically disallowed
                        // Then allow it!
    }

}
